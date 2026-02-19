package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	statikfs "github.com/rakyll/statik/fs"
)

const (
	openAPISpecPath       = "/openapi.json"
	initializerScriptPath = "/swagger-initializer.js"
)

const swaggerInitializerScript = `window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "./openapi.json",
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

func NewHandler(openAPIFilePath string) (http.Handler, error) {
	openAPISpec, err := os.ReadFile(openAPIFilePath)
	if err != nil {
		return nil, fmt.Errorf("read openapi spec %s: %w", openAPIFilePath, err)
	}

	staticFS, err := statikfs.New()
	if err != nil {
		return nil, fmt.Errorf("init swagger-ui static fs: %w", err)
	}

	staticHandler := http.FileServer(staticFS)

	handler := http.NewServeMux()
	handler.HandleFunc(openAPISpecPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		spec, err := prepareOpenAPISpec(openAPISpec, r)
		if err != nil {
			http.Error(w, "failed to prepare openapi spec", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(spec)
	})
	handler.HandleFunc(initializerScriptPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = w.Write([]byte(swaggerInitializerScript))
	})
	handler.Handle("/", staticHandler)

	return handler, nil
}

func prepareOpenAPISpec(raw []byte, r *http.Request) ([]byte, error) {
	var spec map[string]any
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, err
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	spec["host"] = r.Host
	spec["schemes"] = []string{scheme}

	prepared, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}

	return prepared, nil
}
