package user

import (
	"context"
	"log"

	"github.com/fatih/color"

	userv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/models"
)

func (i *Implementation) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	log.Printf("%s: %s: %v, %s: %v, %s: %v, %s: %v, %s: %v",
		color.New(color.FgCyan).Sprint("Create"),
		color.New(color.FgGreen).Sprint("name"), req.GetName(),
		color.New(color.FgGreen).Sprint("email"), req.GetEmail(),
		color.New(color.FgGreen).Sprint("password"), req.GetPassword(),
		color.New(color.FgGreen).Sprint("password_confirm"), req.GetPasswordConfirm(),
		color.New(color.FgGreen).Sprint("role"), req.GetRole(),
	)

	userID, err := i.userService.Create(
		ctx,
		&models.User{
			Name:  req.GetName(),
			Email: req.GetEmail(),
			Role:  models.Role(req.GetRole()),
		},
		req.Password,
		req.PasswordConfirm,
	)
	if err != nil {
		return nil, err
	}

	return &userv1.CreateResponse{Id: userID}, nil
}
