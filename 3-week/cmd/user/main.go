package main

import (
	"context"
	"flag"
	"log"

	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/app/user"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	a := user.NewApp(ctx, configPath)

	if err := a.RunServer(); err != nil {
		log.Fatalf("failed to run gRPC server: %v", err)
	}
}
