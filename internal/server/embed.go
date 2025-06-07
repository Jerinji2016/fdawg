package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

// Embed entire web directory recursively to avoid manual file listing
// This approach automatically includes all files in the web directory,
// eliminating the need to manually add each new template or static file.
// Any new .html templates or .js/.css files will be automatically included.
//
//go:embed web
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
