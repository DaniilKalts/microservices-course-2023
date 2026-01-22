package user

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/user/v1"
	userAPI "github.com/DaniilKalts/microservices-course-2023/3-week/internal/api/grpc/user"
)

type App interface {
	InitDeps(ctx context.Context, configPaths string)

	RunServer() error
}

type app struct {
	serviceProvider ServiceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context, configPath string) App {
	a := &app{}

	a.InitDeps(ctx, configPath)

	return a
}

func (a *app) InitDeps(ctx context.Context, configPaths string) {
	a.serviceProvider = NewServiceProvider(configPaths)

	a.grpcServer = grpc.NewServer()
	reflection.Register(a.grpcServer)
	userv1.RegisterUserV1Server(a.grpcServer, userAPI.NewImplementation(a.serviceProvider.GetUserService(ctx)))
}

func (a *app) RunServer() error {
	defer a.serviceProvider.Close()

	addr := a.serviceProvider.GetConfig().GRPC().Address()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	coloredAddr := color.New(color.FgRed).Sprint(addr)
	log.Printf("gRPC user server is listening on: %s", coloredAddr)

	done := make(chan error, 1)
	go func() {
		if err := a.grpcServer.Serve(lis); err != nil {
			done <- err
		}
		close(done)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-done:
		return err
	case sign := <-stop:
		log.Printf("stopping gRPC server, signal: %v", sign)
		a.grpcServer.GracefulStop()
	}

	return nil
}