package auth

import (
	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type Handler struct {
	authv1.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewHandler(authService service.AuthService) *Handler {
	return &Handler{authService: authService}
}
