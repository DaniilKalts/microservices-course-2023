package operations

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
)

type LoginInput struct {
	Email    string
	Password string
}

func Login(_ context.Context, _ LoginInput) (domainAuth.TokenPair, error) {
	return domainAuth.TokenPair{}, errLoginNotImplemented
}
