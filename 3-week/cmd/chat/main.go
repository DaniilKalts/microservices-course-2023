package main

import (
	"flag"
	"log"
	"net"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	chatv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/chat/v1"
	chatAPI "github.com/DaniilKalts/microservices-course-2023/3-week/internal/api/grpc/chat"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config/env"
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

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen on gRPC chat server: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chatv1.RegisterChatV1Server(s, chatAPI.NewImplementation())

	addr := color.New(color.FgRed).Sprint(grpcConfig.Address())
	log.Printf("gRPC chat server is listening on: %s", addr)

	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve gRPC chat server: %v", err)
	}
}
