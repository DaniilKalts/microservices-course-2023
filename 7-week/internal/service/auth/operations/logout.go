package operations

import (
	"context"
	"fmt"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type LogoutInput struct {
	RefreshToken string
}

func Logout(_ context.Context, jwtManager jwt.Manager, input LogoutInput) error {
	if _, err := jwtManager.VerifyRefreshToken(input.RefreshToken); err != nil {
		return fmt.Errorf("%w: %v", domainAuth.ErrInvalidRefreshToken, err)
	}

	return nil
}
