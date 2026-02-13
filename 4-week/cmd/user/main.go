package main

import (
	"context"
	"flag"
	"log"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/app/user"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	app := user.NewApp(ctx, configPath)

	if err := app.RunServer(); err != nil {
		log.Fatalf("failed to run gRPC server: %v", err)
	}
}
