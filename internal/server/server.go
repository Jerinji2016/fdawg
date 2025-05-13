package server

import (
    "html/template"
    "net/http"
    "path/filepath"

    "github.com/Jerinji2016/fdawg/pkg/utils"
)

// Start initializes and starts the HTTP server
func Start(port string) error {
    // Set up routes
    http.HandleFunc("/", handleIndex)
    
    // Serve static files
    fs := http.FileServer(http.Dir(filepath.Join(utils.ProjectRoot(), "web/static")))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    utils.Info("Server starting on port %s...", port)
    utils.Info("Access the web interface at http://localhost:%s", port)
    return http.ListenAndServe(":"+port, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    // If path is not root, return 404
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    
    // Get the template path
    templatePath := filepath.Join(utils.ProjectRoot(), "web/templates/index.html")
    
    // Parse and execute the template
    tmpl, err := template.ParseFiles(templatePath)
    if err != nil {
        utils.Error("Failed to parse template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    
    err = tmpl.Execute(w, nil)
    if err != nil {
        utils.Error("Failed to execute template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
