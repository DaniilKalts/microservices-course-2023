package chat

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fatih/color"

	chatv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/chat/v1"
)

func (i *Implementation) Create(ctx context.Context, req *chatv1.CreateRequest) (*chatv1.CreateResponse, error) {
	log.Printf("%s: %s: %v", color.New(color.FgCyan).
	Sprint("Create"), color.New(color.FgGreen).
	Sprint("usernames"), req.GetUsernames())

	return &chatv1.CreateResponse{Id: gofakeit.ID()}, nil
}
