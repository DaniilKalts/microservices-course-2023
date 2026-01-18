package chat

import (
	"context"
	"log"

	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/emptypb"

	chatv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/chat/v1"
)

func (i *Implementation) Delete(ctx context.Context, req *chatv1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v", color.New(color.FgCyan).
		Sprint("Delete"), color.New(color.FgGreen).
		Sprint("id"), req.GetId())

	return nil, nil
}
