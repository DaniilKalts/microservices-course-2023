package swaggerui

import "embed"

//go:embed *.html *.css *.js *.png
var StaticFiles embed.FS
