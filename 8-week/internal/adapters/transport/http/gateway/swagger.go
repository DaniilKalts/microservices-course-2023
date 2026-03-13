package gateway

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/http/swagger"
)

type swaggerRoute struct {
	name       string
	basePath   string
	openAPIURL string
}

var swaggerRoutes = []swaggerRoute{
	{
		name:       "merged",
		basePath:   swaggerBasePath,
		openAPIURL: "api/gen/openapi/gateway.swagger.json",
	},
	{
		name:       "user",
		basePath:   swaggerBasePath + "/user",
		openAPIURL: "api/gen/openapi/user/v1/user.swagger.json",
	},
	{
		name:       "profile",
		basePath:   swaggerBasePath + "/profile",
		openAPIURL: "api/gen/openapi/user/v1/profile.swagger.json",
	},
	{
		name:       "auth",
		basePath:   swaggerBasePath + "/auth",
		openAPIURL: "api/gen/openapi/auth/v1/auth.swagger.json",
	},
}

func registerSwaggerHandlers(mux *http.ServeMux) error {
	for _, route := range swaggerRoutes {
		if err := registerSwaggerHandler(mux, route); err != nil {
			return err
		}
	}

	return nil
}

func registerSwaggerHandler(mux *http.ServeMux, route swaggerRoute) error {
	openAPISpec, err := os.ReadFile(route.openAPIURL)
	if err != nil {
		return fmt.Errorf("read %s openapi spec: %w", route.name, err)
	}

	handler := swagger.NewHandler(openAPISpec)

	redirectPath := route.basePath + "/"
	mux.Handle(redirectPath, http.StripPrefix(route.basePath, handler))
	mux.Handle(route.basePath, http.RedirectHandler(redirectPath, http.StatusMovedPermanently))

	return nil
}
