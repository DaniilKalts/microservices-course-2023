package swagger

import (
	"fmt"
	"net/http"

	swaggerui "github.com/DaniilKalts/microservices-course-2023/8-week/web/swagger-ui"
)

const (
	openAPISpecPath       = "/openapi.json"
	initializerScriptPath = "/swagger-initializer.js"
)

const uiInitializerScriptTmpl = `window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: ".%s",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};
`

func NewHandler(openAPISpec []byte) http.Handler {
	uiInitializerScript := []byte(fmt.Sprintf(uiInitializerScriptTmpl, openAPISpecPath))

	mux := http.NewServeMux()
	mux.HandleFunc(openAPISpecPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(openAPISpec)
	})
	mux.HandleFunc(initializerScriptPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = w.Write(uiInitializerScript)
	})
	mux.Handle("/", http.FileServer(http.FS(swaggerui.StaticFiles)))

	return mux
}
