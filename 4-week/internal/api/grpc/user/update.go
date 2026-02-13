package user

import (
	"context"
	"log"

	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/4-week/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
)

func (i *Implementation) Update(ctx context.Context, req *userv1.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("%s: %s: %v, %s: %v, %s: %v",
		color.New(color.FgCyan).Sprint("Update"),
		color.New(color.FgGreen).Sprint("id"), req.GetId(),
		color.New(color.FgGreen).Sprint("name"), req.GetName().GetValue(),
		color.New(color.FgGreen).Sprint("email"), req.GetEmail().GetValue(),
	)

	err := i.userService.Update(
		ctx,
		req.GetId(),
		&models.UpdateUserPatch{
			Name:  &req.GetName().Value,
			Email: &req.GetEmail().Value,
		},
	)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
