package interceptor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	userv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

const (
	authorizationMetadataKey = "authorization"

	userCollectionPath = "/api/v1/users"
	userItemPathPrefix = "/api/v1/users/"
)

var adminOnlyMethods = map[string]struct{}{
	userv1.UserV1_Create_FullMethodName: {},
	userv1.UserV1_List_FullMethodName:   {},
	userv1.UserV1_Get_FullMethodName:    {},
	userv1.UserV1_Update_FullMethodName: {},
	userv1.UserV1_Delete_FullMethodName: {},
}

func AuthInterceptor(jwtManager jwt.Manager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if !requiresAdminRole(info.FullMethod) {
			return handler(ctx, req)
		}

		token, err := accessTokenFromContext(ctx)
		if err != nil {
			return nil, err
		}

		if err = authorizeAdminAccess(token, jwtManager); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func AuthMiddleware(jwtManager jwt.Manager) runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			fullMethod := gatewayMethod(r.Method, r.URL.Path)
			if !requiresAdminRole(fullMethod) {
				next(w, r, pathParams)
				return
			}

			if err := authorizeAdminAccess(r.Header.Get("Authorization"), jwtManager); err != nil {
				writeGatewayError(w, err)
				return
			}

			next(w, r, pathParams)
		}
	}
}

func requiresAdminRole(fullMethod string) bool {
	_, ok := adminOnlyMethods[fullMethod]

	return ok
}

func accessTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "authorization token is required")
	}

	values := md.Get(authorizationMetadataKey)
	for _, value := range values {
		if token := strings.TrimSpace(value); token != "" {
			return token, nil
		}
	}

	return "", status.Error(codes.Unauthenticated, "authorization token is required")
}

func authorizeAdminAccess(token string, jwtManager jwt.Manager) error {
	if jwtManager == nil {
		return status.Error(codes.Internal, "jwt manager is not configured")
	}

	if strings.TrimSpace(token) == "" {
		return status.Error(codes.Unauthenticated, "authorization token is required")
	}

	claims, err := jwtManager.VerifyAccessToken(token)
	if err != nil || claims == nil {
		return status.Error(codes.Unauthenticated, "invalid access token")
	}

	if claims.RoleID != int32(domainUser.RoleAdmin) {
		return status.Error(codes.PermissionDenied, "admin role is required")
	}

	return nil
}

func gatewayMethod(httpMethod, path string) string {
	switch {
	case path == userCollectionPath && httpMethod == http.MethodPost:
		return userv1.UserV1_Create_FullMethodName
	case path == userCollectionPath && httpMethod == http.MethodGet:
		return userv1.UserV1_List_FullMethodName
	case strings.HasPrefix(path, userItemPathPrefix):
		switch httpMethod {
		case http.MethodGet:
			return userv1.UserV1_Get_FullMethodName
		case http.MethodPatch:
			return userv1.UserV1_Update_FullMethodName
		case http.MethodDelete:
			return userv1.UserV1_Delete_FullMethodName
		}
	}

	return ""
}

func writeGatewayError(w http.ResponseWriter, err error) {
	st := status.Convert(err)

	w.Header().Set("Content-Type", "application/json")
	if st.Code() == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", st.Message())
	}

	w.WriteHeader(runtime.HTTPStatusFromCode(st.Code()))
	_, _ = fmt.Fprintf(w, `{"code":%d,"message":%q}`+"\n", st.Code(), st.Message())
}
