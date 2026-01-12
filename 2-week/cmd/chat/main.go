package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	chatv1 "github.com/DaniilKalts/microservices-course-2023/1-week/gen/go/chat/v1"
	"github.com/DaniilKalts/microservices-course-2023/2-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/2-week/internal/config/env"
)

const (
	grpcHostEnvName = "CHAT_GRPC_HOST"
	grpcPortEnvName = "CHAT_GRPC_PORT"
)

type server struct {
	chatv1.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *chatv1.CreateRequest) (*chatv1.CreateResponse, error) {
	log.Printf("%s: %s: %v", color.New(color.FgCyan).Sprint("Create"), color.New(color.FgGreen).Sprint("usernames"), req.GetUsernames())

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &chatv1.CreateResponse{Id: id.String()}, nil
}

func (s *server) Delete(ctx context.Context, req *chatv1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v", color.New(color.FgCyan).Sprint("Delete"), color.New(color.FgGreen).Sprint("id"), req.GetId())
	return nil, nil
}

func (s *server) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v, %s: %v, %s: %v",
		color.New(color.FgCyan).Sprint("SendMessage"),
		color.New(color.FgGreen).Sprint("from"), req.GetFrom(),
		color.New(color.FgGreen).Sprint("text"), req.GetText(),
		color.New(color.FgGreen).Sprint("timestamp"), req.GetTimestamp(),
	)
	return nil, nil
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	flag.Parse()

	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig(grpcHostEnvName, grpcPortEnvName)
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen on gRPC chat server: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chatv1.RegisterChatV1Server(s, &server{})

	addr := color.New(color.FgRed).Sprint(grpcConfig.Address())
	log.Printf("gRPC chat server is listening on: %s", addr)

	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve gRPC chat server: %v", err)
	}
}
