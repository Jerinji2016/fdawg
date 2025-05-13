package server

import (
    "fmt"
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
    return http.ListenAndServe(":"+port, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Flutter Manager Server is running!")
}
