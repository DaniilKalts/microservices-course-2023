package chat

import (
	"context"
	"log"

	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/emptypb"

	chatv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/chat/v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v, %s: %v, %s: %v",
		color.New(color.FgCyan).Sprint("SendMessage"),
		color.New(color.FgGreen).Sprint("from"), req.GetFrom(),
		color.New(color.FgGreen).Sprint("text"), req.GetText(),
		color.New(color.FgGreen).Sprint("timestamp"), req.GetTimestamp(),
	)
	return nil, nil
}