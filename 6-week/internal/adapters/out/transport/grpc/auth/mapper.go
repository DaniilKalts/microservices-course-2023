package auth

import (
	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
)

func toProtoTokenPair(tokens domainAuth.TokenPair) *authv1.TokenPair {
	return &authv1.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
