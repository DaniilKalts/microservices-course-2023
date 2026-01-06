package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/1-week/gen/go/user/v1"
)

const grpcPort = 50052

type server struct {
	userv1.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	log.Printf("%s: %s: %v, %s: %v, %s: %v, %s: %v, %s: %v",
		color.New(color.FgCyan).Sprint("Create"),
		color.New(color.FgGreen).Sprint("name"), req.GetName(),
		color.New(color.FgGreen).Sprint("email"), req.GetEmail(),
		color.New(color.FgGreen).Sprint("password"), req.GetPassword(),
		color.New(color.FgGreen).Sprint("password_confirm"), req.GetPasswordConfirm(),
		color.New(color.FgGreen).Sprint("role"), req.GetRole(),
	)

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &userv1.CreateResponse{Id: id.String()}, nil
}

func (s *server) Get(ctx context.Context, req *userv1.GetRequest) (*userv1.GetResponse, error) {
	log.Printf("%s: %s: %v", color.New(color.FgCyan).Sprint("Get"), color.New(color.FgGreen).Sprint("id"), req.GetId())

	return &userv1.GetResponse{
		User: &userv1.User{
			Id:        gofakeit.ID(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Role:      userv1.Role_USER,
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		}}, nil
}

func (s *server) Update(ctx context.Context, req *userv1.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v, %s: %v, %s: %v",
		color.New(color.FgCyan).Sprint("Update"),
		color.New(color.FgGreen).Sprint("id"), req.GetId(),
		color.New(color.FgGreen).Sprint("name"), req.GetName(),
		color.New(color.FgGreen).Sprint("email"), req.GetEmail(),
	)
	return nil, nil
}

func (s *server) Delete(ctx context.Context, req *userv1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v", color.New(color.FgCyan).Sprint("Delete"), color.New(color.FgGreen).Sprint("id"), req.GetId())
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on gRPC user server: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	userv1.RegisterUserV1Server(s, &server{})

	addr := color.New(color.FgRed).Sprintf("localhost:%d", grpcPort)
	log.Printf("gRPC user server is listening on: %s", addr)

	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve gRPC user server: %v", err)
	}
}
