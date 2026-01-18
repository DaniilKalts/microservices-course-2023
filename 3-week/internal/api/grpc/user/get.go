package user

import (
	"context"
	"log"

	"github.com/fatih/color"

	userv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/converter"
)

func (i *Implementation) Get(ctx context.Context, req *userv1.GetRequest) (*userv1.GetResponse, error) {
	log.Printf("%s: %s: %v",
		color.New(color.FgCyan).Sprint("Get"),
		color.New(color.FgGreen).Sprint("id"),
		req.GetId(),
	)

	user, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userv1.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
