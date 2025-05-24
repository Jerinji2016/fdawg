package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/Jerinji2016/fdawg/internal/server/api"
	"github.com/Jerinji2016/fdawg/pkg/environment"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/utils"
)

// ServerData contains data to be passed to templates
type ServerData struct {
	Project         *flutter.ValidationResult
	ActivePage      string
	EnvFiles        []environment.EnvFile
	SelectedEnvFile *environment.EnvFile
}

// Start initializes and starts the HTTP server
func Start(port string, project *flutter.ValidationResult) error {
	// Create base server data
	baseData := &ServerData{
		Project: project,
	}

	// Set up API routes
	api.SetupAPIRoutes(project)

	// Set up routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := *baseData
		data.ActivePage = "overview"
		handlePage(w, r, &data, "overview")
	})

	http.HandleFunc("/environment", func(w http.ResponseWriter, r *http.Request) {
		data := *baseData
		data.ActivePage = "environment"

		// Get environment files
		envFiles, err := environment.ListEnvFiles(project.ProjectPath)
		if err != nil {
			utils.Error("Failed to list environment files: %v", err)
		} else {
			data.EnvFiles = envFiles

			// Set selected environment file if available
			if len(envFiles) > 0 {
				// Check if there's a query parameter for selected environment
				selectedEnv := r.URL.Query().Get("env")
				if selectedEnv != "" {
					// Find the selected environment file
					for _, envFile := range envFiles {
						if envFile.Name == selectedEnv {
							data.SelectedEnvFile = &envFile
							break
						}
					}
				}

				// If no environment is selected, use the first one
				if data.SelectedEnvFile == nil {
					data.SelectedEnvFile = &envFiles[0]
				}
			}
		}

		handlePage(w, r, &data, "environment")
	})

	http.HandleFunc("/assets", func(w http.ResponseWriter, r *http.Request) {
		data := *baseData
		data.ActivePage = "assets"
		handlePage(w, r, &data, "assets")
	})

	http.HandleFunc("/localizations", func(w http.ResponseWriter, r *http.Request) {
		data := *baseData
		data.ActivePage = "localizations"
		handlePage(w, r, &data, "localizations")
	})

	http.HandleFunc("/fastlane", func(w http.ResponseWriter, r *http.Request) {
		data := *baseData
		data.ActivePage = "fastlane"
		handlePage(w, r, &data, "fastlane")
	})

	http.HandleFunc("/run-configs", func(w http.ResponseWriter, r *http.Request) {
		data := *baseData
		data.ActivePage = "run-configs"
		handlePage(w, r, &data, "run_configs")
	})

	http.HandleFunc("/namer", func(w http.ResponseWriter, r *http.Request) {
		data := *baseData
		data.ActivePage = "namer"
		handlePage(w, r, &data, "namer")
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

// handlePage renders a page using the layout template and the specified content template
func handlePage(w http.ResponseWriter, r *http.Request, data *ServerData, templateName string) {
	tmplFS := GetTemplateFS()

	// Read the layout template
	layoutData, err := fs.ReadFile(tmplFS, "layout.html")
	if err != nil {
		utils.Error("Failed to read layout template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Read the content template
	contentData, err := fs.ReadFile(tmplFS, templateName+".html")
	if err != nil {
		utils.Error("Failed to read %s template: %v", templateName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a new template with the layout
	tmpl, err := template.New("layout").Parse(string(layoutData))
	if err != nil {
		utils.Error("Failed to parse layout template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse the content template
	_, err = tmpl.New("content").Parse(string(contentData))
	if err != nil {
		utils.Error("Failed to parse %s template: %v", templateName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, data)
	if err != nil {
		utils.Error("Failed to execute template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
