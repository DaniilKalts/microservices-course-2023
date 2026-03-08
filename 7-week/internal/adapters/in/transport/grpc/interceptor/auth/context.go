package auth

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type claimsContextKey struct{}

func ClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*jwt.Claims)
	return claims, ok
}

func withClaims(ctx context.Context, claims *jwt.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}
