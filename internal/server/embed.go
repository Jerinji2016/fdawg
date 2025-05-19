package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed web/templates/layout.html
//go:embed web/templates/overview.html
//go:embed web/templates/environment.html
//go:embed web/templates/assets.html
//go:embed web/templates/localizations.html
//go:embed web/templates/fastlane.html
//go:embed web/templates/run_configs.html
//go:embed web/static/css/style.css
//go:embed web/static/js/main.js
//go:embed web/static/js/environment.js
//go:embed web/static/js/assets.js
//go:embed web/static/js/toast.js
var webFS embed.FS

// GetTemplateFS returns a filesystem for accessing embedded templates
func GetTemplateFS() fs.FS {
	templates, err := fs.Sub(webFS, "web/templates")
	if err != nil {
		log.Printf("Error creating template filesystem: %v", err)
		return nil
	}
	return templates
}

// GetStaticFS returns a filesystem for serving embedded static files
func GetStaticFS() http.FileSystem {
	static, err := fs.Sub(webFS, "web/static")
	if err != nil {
		log.Printf("Error creating static filesystem: %v", err)
		return nil
	}
	return http.FS(static)
}
