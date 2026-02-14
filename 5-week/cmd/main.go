package main

import (
	"context"
	"flag"
	"log"

	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	a, err := app.New(ctx, configPath)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err = a.Run(ctx); err != nil {
		log.Fatalf("failed to run gRPC server: %v", err)
	}	
}
