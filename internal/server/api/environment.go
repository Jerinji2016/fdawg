package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/environment"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
)

// EnvironmentAPI handles environment-related API endpoints
type EnvironmentAPI struct {
	project *flutter.ValidationResult
}

// NewEnvironmentAPI creates a new EnvironmentAPI instance
func NewEnvironmentAPI(project *flutter.ValidationResult) *EnvironmentAPI {
	return &EnvironmentAPI{
		project: project,
	}
}

// RegisterRoutes registers environment API routes
func (api *EnvironmentAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/environment/create", api.handleCreateEnvironment)
	mux.HandleFunc("/api/environment/add-variable", api.handleAddVariable)
	mux.HandleFunc("/api/environment/delete-variable", api.handleDeleteVariable)
	mux.HandleFunc("/api/environment/delete-env", api.handleDeleteEnvironment)
	mux.HandleFunc("/api/environment/download", api.handleDownloadEnvironment)
}

// handleCreateEnvironment handles POST requests to create environment files
func (api *EnvironmentAPI) handleCreateEnvironment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get environment name
	envName := r.FormValue("env_name")
	if envName == "" {
		http.Error(w, "Environment name is required", http.StatusBadRequest)
		return
	}

	// Get copy from parameter
	copyFrom := r.FormValue("copy_from")

	// Create environment file
	var createErr error
	if copyFrom != "" {
		// Copy from existing environment file
		createErr = environment.CopyEnvFile(api.project.ProjectPath, copyFrom, envName)
	} else {
		// Create empty environment file
		createErr = environment.CreateEnvFile(api.project.ProjectPath, envName, make(map[string]interface{}))
	}

	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed to create environment file: %v", createErr), http.StatusInternalServerError)
		return
	}

	// Redirect back to the environment page
	http.Redirect(w, r, fmt.Sprintf("/environment?env=%s", envName), http.StatusSeeOther)
}

// handleAddVariable handles POST requests to add variables to environment files
func (api *EnvironmentAPI) handleAddVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get environment name
	envName := r.FormValue("env_name")
	if envName == "" {
		http.Error(w, "Environment name is required", http.StatusBadRequest)
		return
	}

	// Get key and value
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	// Validate key format (must start with letter or underscore, and contain only letters, numbers, and underscores)
	keyRegex := regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
	if !keyRegex.MatchString(key) {
		if regexp.MustCompile(`^\d`).MatchString(key) {
			http.Error(w, "Key must not start with a number (Dart variable naming convention)", http.StatusBadRequest)
		} else {
			http.Error(w, "Key must contain only letters, numbers, and underscores (no spaces or special characters)", http.StatusBadRequest)
		}
		return
	}

	valueStr := r.FormValue("value")

	// Parse value (try to convert to appropriate type)
	var value interface{} = valueStr

	// Try to parse as number or boolean
	if strings.EqualFold(valueStr, "true") {
		value = true
	} else if strings.EqualFold(valueStr, "false") {
		value = false
	} else if strings.Contains(valueStr, ".") {
		// Try to parse as float
		var f float64
		if _, err := fmt.Sscanf(valueStr, "%f", &f); err == nil {
			value = f
		}
	} else {
		// Try to parse as integer
		var i int64
		if _, err := fmt.Sscanf(valueStr, "%d", &i); err == nil {
			value = i
		}
	}

	// Add variable to environment file
	err = environment.AddVariable(api.project.ProjectPath, envName, key, value)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add variable: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect back to the environment page
	http.Redirect(w, r, fmt.Sprintf("/environment?env=%s", envName), http.StatusSeeOther)
}

// handleDeleteVariable handles POST requests to delete variables from environment files
func (api *EnvironmentAPI) handleDeleteVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get environment name
	envName := r.FormValue("env_name")
	if envName == "" {
		http.Error(w, "Environment name is required", http.StatusBadRequest)
		return
	}

	// Get key
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	// Delete variable from environment file
	err = environment.DeleteVariable(api.project.ProjectPath, envName, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete variable: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect back to the environment page
	http.Redirect(w, r, fmt.Sprintf("/environment?env=%s", envName), http.StatusSeeOther)
}

// handleDeleteEnvironment handles POST requests to delete environment files
func (api *EnvironmentAPI) handleDeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get environment name
	envName := r.FormValue("env_name")
	if envName == "" {
		http.Error(w, "Environment name is required", http.StatusBadRequest)
		return
	}

	// Delete environment file
	err = environment.DeleteEnvFile(api.project.ProjectPath, envName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete environment file: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect back to the environment page
	http.Redirect(w, r, "/environment", http.StatusSeeOther)
}

// handleDownloadEnvironment handles GET requests to download environment files
func (api *EnvironmentAPI) handleDownloadEnvironment(w http.ResponseWriter, r *http.Request) {
	// Get environment name
	envName := r.URL.Query().Get("env_name")
	if envName == "" {
		http.Error(w, "Environment name is required", http.StatusBadRequest)
		return
	}

	// Get environment file
	envFile, err := environment.GetEnvFile(api.project.ProjectPath, envName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get environment file: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", envName))
	w.Header().Set("Content-Type", "application/json")

	// Marshal the variables to JSON
	data, err := json.MarshalIndent(envFile.Variables, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the JSON to the response
	w.Write(data)
}

// SetupEnvironmentAPIRoutes sets up environment API routes
func SetupEnvironmentAPIRoutes(project *flutter.ValidationResult) {
	environmentAPI := NewEnvironmentAPI(project)
	environmentAPI.RegisterRoutes(http.DefaultServeMux)
}
