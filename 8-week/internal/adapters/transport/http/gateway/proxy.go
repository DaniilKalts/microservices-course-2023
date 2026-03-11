package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/DaniilKalts/microservices-course-2023/8-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/http/swagger"
	appconfig "github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
)

const swaggerBasePath = "/swagger"

type Config struct {
	GRPCAddress string
	TLS         appconfig.TLSConfig
}

type Proxy struct {
	handler http.Handler
	conn    *grpc.ClientConn
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.handler.ServeHTTP(w, r)
}

func (p *Proxy) Close() error {
	if p == nil || p.conn == nil {
		return nil
	}

	return p.conn.Close()
}

type swaggerRoute struct {
	name       string
	basePath   string
	openAPIURL string
}

var swaggerRoutes = []swaggerRoute{
	{name: "merged", basePath: swaggerBasePath, openAPIURL: "gen/openapi/gateway.swagger.json"},
	{name: "user", basePath: swaggerBasePath + "/user", openAPIURL: "gen/openapi/user/v1/user.swagger.json"},
	{name: "profile", basePath: swaggerBasePath + "/profile", openAPIURL: "gen/openapi/user/v1/profile.swagger.json"},
	{name: "auth", basePath: swaggerBasePath + "/auth", openAPIURL: "gen/openapi/auth/v1/auth.swagger.json"},
}

func NewProxy(ctx context.Context, cfg Config) (*Proxy, error) {
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
	if err := userv1.RegisterProfileV1Handler(ctx, gatewayMux, conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("register profile grpc-gateway handler: %w", err)
	}
	if err := authv1.RegisterAuthV1Handler(ctx, gatewayMux, conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("register auth grpc-gateway handler: %w", err)
	}

	if err := registerSwaggerHandlers(mux); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &Proxy{
		handler: WithTracing(mux, opentracing.GlobalTracer()),
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

func grpcGatewayDialOptions(tlsCfg appconfig.TLSConfig) ([]grpc.DialOption, error) {
	traceInterceptor := grpc.WithChainUnaryInterceptor(TracingClientInterceptor(opentracing.GlobalTracer()))

	if !tlsCfg.Enabled {
		return []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			traceInterceptor,
		}, nil
	}

	clientCreds, err := credentials.NewClientTLSFromFile(tlsCfg.CertFile, "")
	if err != nil {
		return nil, fmt.Errorf("load gateway grpc tls cert: %w", err)
	}

	return []grpc.DialOption{
		grpc.WithTransportCredentials(clientCreds),
		traceInterceptor,
	}, nil
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
	openAPISpec, err := os.ReadFile(route.openAPIURL)
	if err != nil {
		return fmt.Errorf("read %s openapi spec: %w", route.name, err)
	}

	handler, err := swagger.NewHandler(openAPISpec)
	if err != nil {
		return fmt.Errorf("init %s swagger-ui handler: %w", route.name, err)
	}

	redirectPath := route.basePath + "/"
	mux.Handle(redirectPath, http.StripPrefix(route.basePath, handler))
	mux.Handle(route.basePath, http.RedirectHandler(redirectPath, http.StatusMovedPermanently))

	return nil
}
