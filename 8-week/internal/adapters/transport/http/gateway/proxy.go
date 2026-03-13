package gateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	authv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/http/gateway/interceptor"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/http/gateway/middleware"
	appconfig "github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
)

const (
	swaggerBasePath   = "/swagger"
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
)

type Config struct {
	GRPCAddress    string
	GatewayAddress string
	TLS            appconfig.TLSConfig
	Tracer         opentracing.Tracer
	CircuitBreaker interceptor.CircuitBreakerConfig
}

type Proxy struct {
	server *http.Server
	conn   *grpc.ClientConn
	tls    appconfig.TLSConfig
}

func (p *Proxy) Addr() string {
	return p.server.Addr
}

func (p *Proxy) Serve() error {
	var err error
	if p.tls.Enabled {
		err = p.server.ListenAndServeTLS(p.tls.CertFile, p.tls.KeyFile)
	} else {
		err = p.server.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("gateway server: %w", err)
	}

	return nil
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	var errs []error

	if p.server != nil {
		if err := p.server.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown gateway http server: %w", err))
		}
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close gateway grpc conn: %w", err))
		}
	}

	return errors.Join(errs...)
}

func NewProxy(ctx context.Context, cfg Config) (_ *Proxy, err error) {
	tracer := cfg.Tracer
	if tracer == nil {
		tracer = opentracing.NoopTracer{}
	}

	gatewayMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: false,
				},
			},
		}),
	)

	grpcEndpoint, err := grpcGatewayEndpoint(cfg.GRPCAddress)
	if err != nil {
		return nil, fmt.Errorf("prepare grpc endpoint for gateway: %w", err)
	}

	dialOpts, err := grpcGatewayDialOptions(cfg.TLS, tracer, cfg.CircuitBreaker)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(grpcEndpoint, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("dial grpc endpoint for gateway: %w", err)
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/", gatewayMux)

	if err = userv1.RegisterUserV1Handler(ctx, gatewayMux, conn); err != nil {
		return nil, fmt.Errorf("register user grpc-gateway handler: %w", err)
	}
	if err = userv1.RegisterProfileV1Handler(ctx, gatewayMux, conn); err != nil {
		return nil, fmt.Errorf("register profile grpc-gateway handler: %w", err)
	}
	if err = authv1.RegisterAuthV1Handler(ctx, gatewayMux, conn); err != nil {
		return nil, fmt.Errorf("register auth grpc-gateway handler: %w", err)
	}

	if err = registerSwaggerHandlers(mux); err != nil {
		return nil, err
	}

	return &Proxy{
		server: &http.Server{
			Addr:              cfg.GatewayAddress,
			Handler:           middleware.WithTracing(mux, tracer),
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
		},
		conn: conn,
		tls:  cfg.TLS,
	}, nil
}

func grpcGatewayEndpoint(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", err
	}

	if isWildcardHost(host) {
		host = "localhost"
	}

	return net.JoinHostPort(host, port), nil
}

func isWildcardHost(host string) bool {
	return host == "" || host == "0.0.0.0" || host == "::"
}

func grpcGatewayDialOptions(tlsCfg appconfig.TLSConfig, tracer opentracing.Tracer, cbCfg interceptor.CircuitBreakerConfig) ([]grpc.DialOption, error) {
	interceptors := grpc.WithChainUnaryInterceptor(
		interceptor.CircuitBreakerInterceptor(cbCfg),
		interceptor.TracingInterceptor(tracer),
	)

	if !tlsCfg.Enabled {
		return []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			interceptors,
		}, nil
	}

	clientCreds, err := credentials.NewClientTLSFromFile(tlsCfg.CertFile, "")
	if err != nil {
		return nil, fmt.Errorf("load gateway grpc tls cert: %w", err)
	}

	return []grpc.DialOption{
		grpc.WithTransportCredentials(clientCreds),
		interceptors,
	}, nil
}
