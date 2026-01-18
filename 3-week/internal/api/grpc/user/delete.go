package user

import (
	"context"
	"log"

	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/user/v1"
)

func (i *Implementation) Delete(ctx context.Context, req *userv1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v",
		color.New(color.FgCyan).Sprint("Delete"),
		color.New(color.FgGreen).Sprint("id"),
		req.GetId(),
	)

	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
