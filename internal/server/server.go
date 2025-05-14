package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/utils"
)

// ServerData contains data to be passed to templates
type ServerData struct {
    Project *flutter.ValidationResult
}

// Start initializes and starts the HTTP server
func Start(port string, project *flutter.ValidationResult) error {
    // Create server data
    data := &ServerData{
        Project: project,
    }

    // Set up routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        handleIndex(w, r, data)
    })

    // Serve static files from embedded filesystem
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(GetStaticFS())))

    // Format the address with the port
    addr := fmt.Sprintf(":%s", port)

    utils.Info("Server starting on port %s...", port)
    utils.Info("Access the web interface at http://localhost:%s", port)
    utils.Info("Press Ctrl+C to stop the server")
    
    return http.ListenAndServe(addr, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request, data *ServerData) {
    // If path is not root, return 404
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    // Parse and execute the template from embedded filesystem
    tmplFS := GetTemplateFS()
    tmplData, err := fs.ReadFile(tmplFS, "index.html")
    if err != nil {
        utils.Error("Failed to read template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.New("index").Parse(string(tmplData))
    if err != nil {
        utils.Error("Failed to parse template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        utils.Error("Failed to execute template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
