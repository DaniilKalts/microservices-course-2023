package grpc

import (
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/interceptor"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type ServerConfig struct {
	EnableTLS bool
	CertFile  string
	KeyFile   string
}

type Deps struct {
	Config     ServerConfig
	JWTManager jwt.Manager
	Logger     *zap.Logger
}

func NewServer(deps Deps) (*grpc.Server, error) {
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingInterceptor(deps.Logger),
			interceptor.AuthInterceptor(deps.JWTManager),
			interceptor.ValidationInterceptor(),
		),
	}

	cfg := deps.Config

	if cfg.EnableTLS {
		if cfg.CertFile == "" || cfg.KeyFile == "" {
			return nil, errors.New("grpc tls is enabled but cert or key file is empty")
		}

		creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, err
		}

		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	return grpc.NewServer(grpcOpts...), nil
}
