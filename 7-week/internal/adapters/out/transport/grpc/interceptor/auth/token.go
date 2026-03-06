package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
)

const (
	authorizationMetadataKey = "authorization"
)

func accessTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, domainAuth.ErrAuthorizationTokenMissing.Error())
	}

	values := md.Get(authorizationMetadataKey)
	for _, value := range values {
		if token := strings.TrimSpace(value); token != "" {
			return token, nil
		}
	}

	return "", status.Error(codes.Unauthenticated, domainAuth.ErrAuthorizationTokenMissing.Error())
}
