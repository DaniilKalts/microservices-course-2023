package grpc

import (
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/interceptor"
)

type ServerConfig struct {
	EnableTLS bool
	CertFile  string
	KeyFile   string
}

func NewServer(cfg ServerConfig) (*grpc.Server, error) {
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.ValidationInterceptor(),
		),
	}

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
