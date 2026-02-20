package auth

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
)

func (s *service) Login(_ context.Context, _, _ string) (domainAuth.TokenPair, error) {
	return domainAuth.TokenPair{}, errLoginNotImplemented
}
