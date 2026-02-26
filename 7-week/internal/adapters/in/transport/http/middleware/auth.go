package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

const (
	userCollectionPath = "/api/v1/users"
	userItemPathPrefix = "/api/v1/users/"
)

var authenticatedMethods = map[string]struct{}{}

var adminOnlyMethods = map[string]struct{}{
	userv1.UserV1_Create_FullMethodName: {},
	userv1.UserV1_Update_FullMethodName: {},
	userv1.UserV1_Delete_FullMethodName: {},
}

func AuthMiddleware(jwtManager jwt.Manager) runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			fullMethod := gatewayMethod(r.Method, r.URL.Path)
			requiredRole, requiresAuth := requiredRole(fullMethod)
			if !requiresAuth {
				next(w, r, pathParams)
				return
			}

			if err := authorize(r.Header.Get("Authorization"), jwtManager, requiredRole); err != nil {
				writeGatewayError(w, err)
				return
			}

			next(w, r, pathParams)
		}
	}
}

func requiredRole(fullMethod string) (domainUser.Role, bool) {
	if _, ok := authenticatedMethods[fullMethod]; ok {
		return domainUser.RoleUser, true
	}

	if _, ok := adminOnlyMethods[fullMethod]; ok {
		return domainUser.RoleAdmin, true
	}

	return 0, false
}

func authorize(token string, jwtManager jwt.Manager, requiredRole domainUser.Role) error {
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

	if !hasRequiredRole(claims.RoleID, requiredRole) {
		return status.Error(codes.PermissionDenied, "insufficient role permissions")
	}

	return nil
}

func hasRequiredRole(roleID int32, requiredRole domainUser.Role) bool {
	switch requiredRole {
	case domainUser.RoleUser:
		return roleID == int32(domainUser.RoleUser) || roleID == int32(domainUser.RoleAdmin)
	case domainUser.RoleAdmin:
		return roleID == int32(domainUser.RoleAdmin)
	default:
		return false
	}
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
