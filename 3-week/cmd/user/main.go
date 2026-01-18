package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config/env"
)

type server struct {
	userv1.UnimplementedUserV1Server
	pool *pgxpool.Pool
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	builderCreate := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "name", "email", "password_hash", "role").
		Values(id, req.GetName(), req.GetEmail(), hashedPassword, req.GetRole()).
		Suffix("RETURNING id")

	query, args, err := builderCreate.ToSql()
	if err != nil {
		return nil, err
	}

	var userID string
	if err := s.pool.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		return nil, err
	}

	return &userv1.CreateResponse{Id: userID}, nil
}

func (s *server) Get(ctx context.Context, req *userv1.GetRequest) (*userv1.GetResponse, error) {
	log.Printf("%s: %s: %v",
		color.New(color.FgCyan).Sprint("Get"),
		color.New(color.FgGreen).Sprint("id"),
		req.GetId(),
	)

	var user userv1.User
	var createdAt, updatedAt time.Time

	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	err = s.pool.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Role,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	return &userv1.GetResponse{User: &user}, nil
}

func (s *server) Update(ctx context.Context, req *userv1.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v, %s: %v, %s: %v",
		color.New(color.FgCyan).Sprint("Update"),
		color.New(color.FgGreen).Sprint("id"), req.GetId(),
		color.New(color.FgGreen).Sprint("name"), req.GetName().GetValue(),
		color.New(color.FgGreen).Sprint("email"), req.GetEmail().GetValue(),
	)

	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.GetId()})

	if req.GetName() != nil {
		builderUpdate = builderUpdate.Set("name", req.GetName().GetValue())
	}

	if req.GetEmail() != nil {
		builderUpdate = builderUpdate.Set("email", req.GetEmail().GetValue())
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *userv1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v",
		color.New(color.FgCyan).Sprint("Delete"),
		color.New(color.FgGreen).Sprint("id"),
		req.GetId(),
	)

	builderDelete := sq.Delete("users").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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

	s := grpc.NewServer()
	reflection.Register(s)
	userv1.RegisterUserV1Server(s, &server{pool: pool})

	addr := color.New(color.FgRed).Sprint(grpcConfig.Address())
	log.Printf("gRPC user server is listening on: %s", addr)

	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve gRPC user server: %v", err)
	}
}
