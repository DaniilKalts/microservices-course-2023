package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/user/v1"
	userAPI "github.com/DaniilKalts/microservices-course-2023/3-week/internal/api/grpc/user"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config/env"
	userRepository "github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository/user"
	userService "github.com/DaniilKalts/microservices-course-2023/3-week/internal/service/user"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	flag.Parse()

	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen on gRPC user server: %v", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, postgresConfig.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	userRepo := userRepository.NewRepository(pool)
	userService := userService.NewService(userRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	userv1.RegisterUserV1Server(s, userAPI.NewImplementation(userService))

	addr := color.New(color.FgRed).Sprint(grpcConfig.Address())
	log.Printf("gRPC user server is listening on: %s", addr)

	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve gRPC user server: %v", err)
	}
}
