package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationMetadataKey = "authorization"
)

func accessTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, ErrAuthorizationTokenMissing.Error())
	}

	values := md.Get(authorizationMetadataKey)
	for _, value := range values {
		token := extractBearerToken(value)
		if token != "" {
			return token, nil
		}
	}

	return "", status.Error(codes.Unauthenticated, ErrAuthorizationTokenMissing.Error())
}

func extractBearerToken(value string) string {
	value = strings.TrimSpace(value)

	scheme, token, found := strings.Cut(value, " ")
	if found && strings.EqualFold(scheme, "Bearer") {
		return strings.TrimSpace(token)
	}

	if strings.EqualFold(value, "Bearer") {
		return ""
	}

	return value
}
