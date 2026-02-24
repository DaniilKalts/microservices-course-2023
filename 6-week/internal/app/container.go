package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/in/database/postgres"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/in/transport/http/swagger"
	grpcTransport "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc"
	authAPI "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/handlers/auth"
	userAPI "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/handlers/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database/transaction"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config/env"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

const swaggerBasePath = "/swagger"

type swaggerRoute struct {
	name       string
	basePath   string
	openAPIURL string
}

var swaggerRoutes = []swaggerRoute{
	{name: "merged", basePath: swaggerBasePath, openAPIURL: "gen/openapi/gateway.swagger.json"},
	{name: "user", basePath: swaggerBasePath + "/user", openAPIURL: "gen/openapi/user/v1/user.swagger.json"},
	{name: "auth", basePath: swaggerBasePath + "/auth", openAPIURL: "gen/openapi/auth/v1/auth.swagger.json"},
}

type Container struct {
	Cfg config.Config

	DB         database.Client
	Tx         database.TxManager
	JWTManager jwt.Manager

	Repositories repository.Repositories
	Services     service.Services

	userHandler userv1.UserV1Server

	authHandler authv1.AuthV1Server

	GRPC    *grpc.Server
	Gateway http.Handler
}

func Build(ctx context.Context, configPath string) (*Container, error) {
	container := &Container{}

	if err := container.initConfig(configPath); err != nil {
		return nil, err
	}
	if err := container.initDatabase(ctx); err != nil {
		return nil, err
	}
	if err := container.initJWTManager(); err != nil {
		return nil, err
	}

	container.initTxManager()
	container.initRepositories()
	container.initServices()
	container.initHandlers()

	if err := container.initGRPC(); err != nil {
		return nil, err
	}
	if err := container.initGateway(ctx); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initConfig(configPath string) error {
	if err := config.Load(configPath); err != nil {
		return fmt.Errorf("load dotenv config: %w", err)
	}

	cfg, err := env.NewConfig()
	if err != nil {
		return fmt.Errorf("load env config: %w", err)
	}

	c.Cfg = cfg

	return nil
}

func (c *Container) initDatabase(ctx context.Context) error {
	db, err := postgres.New(ctx, c.Cfg.Postgres().DSN())
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}

	c.DB = db

	return nil
}

func (c *Container) initTxManager() {
	c.Tx = transaction.NewTransactionManager(c.DB.DB())
}

func (c *Container) initJWTManager() error {
	privateKey, err := jwt.LoadPrivateKey(c.Cfg.JWT().PrivateKeyFile())
	if err != nil {
		return fmt.Errorf("load jwt private key: %w", err)
	}

	publicKey, err := jwt.LoadPublicKey(c.Cfg.JWT().PublicKeyFile())
	if err != nil {
		return fmt.Errorf("load jwt public key: %w", err)
	}

	jwtManager, err := jwt.NewManager(privateKey, publicKey, jwt.Config{
		Issuer:          c.Cfg.JWT().Issuer(),
		Subject:         c.Cfg.JWT().Subject(),
		Audience:        c.Cfg.JWT().Audience(),
		AccessTokenTTL:  c.Cfg.JWT().AccessExpiresAt(),
		RefreshTokenTTL: c.Cfg.JWT().RefreshExpiresAt(),
		NotBeforeOffset: c.Cfg.JWT().NotBefore(),
		IssuedAtOffset:  c.Cfg.JWT().IssuedAt(),
	})
	if err != nil {
		return fmt.Errorf("init jwt manager: %w", err)
	}

	c.JWTManager = jwtManager

	return nil
}

func (c *Container) initRepositories() {
	c.Repositories = repository.NewRepositories(repository.Deps{DB: c.DB})
}

func (c *Container) initServices() {
	c.Services = service.NewServices(service.Deps{
		Repositories: c.Repositories,
		JWTManager:   c.JWTManager,
	})
}

func (c *Container) initHandlers() {
	c.userHandler = userAPI.NewHandler(c.Services.User)
	c.authHandler = authAPI.NewHandler(c.Services.Auth)
}

func (c *Container) initGRPC() error {
	grpcServer, err := grpcTransport.NewServer(grpcTransport.ServerConfig{
		EnableTLS: c.Cfg.TLS().Enabled(),
		CertFile:  c.Cfg.TLS().CertFile(),
		KeyFile:   c.Cfg.TLS().KeyFile(),
	})
	if err != nil {
		return err
	}

	grpcTransport.RegisterServices(grpcServer, grpcTransport.Handlers{
		User: c.userHandler,
		Auth: c.authHandler,
	})

	c.GRPC = grpcServer

	return nil
}

func (c *Container) initGateway(ctx context.Context) error {
	gatewayMux := runtime.NewServeMux()

	mux := http.NewServeMux()
	mux.Handle("/", gatewayMux)

	if err := userv1.RegisterUserV1HandlerServer(ctx, gatewayMux, c.userHandler); err != nil {
		return fmt.Errorf("register user grpc-gateway handler: %w", err)
	}
	if err := authv1.RegisterAuthV1HandlerServer(ctx, gatewayMux, c.authHandler); err != nil {
		return fmt.Errorf("register auth grpc-gateway handler: %w", err)
	}

	if err := registerSwaggerHandlers(mux); err != nil {
		return err
	}

	c.Gateway = mux

	return nil
}

func registerSwaggerHandlers(mux *http.ServeMux) error {
	for _, route := range swaggerRoutes {
		if err := registerSwaggerHandler(mux, route); err != nil {
			return err
		}
	}

	return nil
}

func registerSwaggerHandler(mux *http.ServeMux, route swaggerRoute) error {
	handler, err := swagger.NewHandler(route.openAPIURL)
	if err != nil {
		return fmt.Errorf("init %s swagger-ui handler: %w", route.name, err)
	}

	redirectPath := route.basePath + "/"
	mux.Handle(redirectPath, http.StripPrefix(route.basePath, handler))
	mux.Handle(route.basePath, http.RedirectHandler(redirectPath, http.StatusMovedPermanently))

	return nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
