package gateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/transport/http/swagger"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
)

const swaggerBasePath = "/swagger"

type Config struct {
	GRPCAddress string
	TLS         config.TLSConfig
}

type proxyHandler struct {
	http.Handler
	conn *grpc.ClientConn
}

func (h *proxyHandler) Close() error {
	if h == nil || h.conn == nil {
		return nil
	}

	return h.conn.Close()
}

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

func NewProxy(ctx context.Context, cfg Config) (http.Handler, error) {
	gatewayMux := runtime.NewServeMux()

	grpcEndpoint, err := grpcGatewayEndpoint(cfg.GRPCAddress)
	if err != nil {
		return nil, fmt.Errorf("prepare grpc endpoint for gateway: %w", err)
	}

	dialOpts, err := grpcGatewayDialOptions(cfg.TLS)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(grpcEndpoint, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("dial grpc endpoint for gateway: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gatewayMux)

	if err := userv1.RegisterUserV1Handler(ctx, gatewayMux, conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("register user grpc-gateway handler: %w", err)
	}
	if err := authv1.RegisterAuthV1Handler(ctx, gatewayMux, conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("register auth grpc-gateway handler: %w", err)
	}

	if err := registerSwaggerHandlers(mux); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &proxyHandler{
		Handler: mux,
		conn:    conn,
	}, nil
}

func grpcGatewayEndpoint(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", err
	}

	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "localhost"
	}

	return net.JoinHostPort(host, port), nil
}

func grpcGatewayDialOptions(tlsCfg config.TLSConfig) ([]grpc.DialOption, error) {
	if tlsCfg == nil {
		return nil, errors.New("tls config is nil")
	}

	if !tlsCfg.Enabled() {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}

	clientCreds, err := credentials.NewClientTLSFromFile(tlsCfg.CertFile(), "")
	if err != nil {
		return nil, fmt.Errorf("load gateway grpc tls cert: %w", err)
	}

	return []grpc.DialOption{grpc.WithTransportCredentials(clientCreds)}, nil
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
