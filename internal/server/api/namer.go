package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Jerinji2016/fdawg/internal/server/helpers"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/namer"
)

// NamerAPI handles namer-related API endpoints
type NamerAPI struct{}

// NewNamerAPI creates a new NamerAPI instance
func NewNamerAPI() *NamerAPI {
	return &NamerAPI{}
}

// RegisterRoutes registers namer API routes
func (api *NamerAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/namer/get", api.handleGetAppNames)
	mux.HandleFunc("/api/namer/set", api.handleSetAppNames)
	mux.HandleFunc("/api/namer/platforms", api.handleGetPlatforms)
}

// GetAppNamesRequest represents a request to get app names
type GetAppNamesRequest struct {
	Platforms []string `json:"platforms,omitempty"`
}

// SetAppNamesAPIRequest represents a request to set app names via API
type SetAppNamesAPIRequest struct {
	Universal string            `json:"universal,omitempty"`
	Platforms map[string]string `json:"platforms,omitempty"`
}

// PlatformsResponse represents available platforms response
type PlatformsResponse struct {
	Available []PlatformInfo `json:"available"`
	All       []PlatformInfo `json:"all"`
}

// PlatformInfo represents platform information
type PlatformInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Available   bool   `json:"available"`
	Description string `json:"description"`
}

// handleGetAppNames handles GET requests to retrieve app names
func (api *NamerAPI) handleGetAppNames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate Flutter project
	result, err := flutter.ValidateProject(".")
	if err != nil || !result.IsValid {
		http.Error(w, "Not a valid Flutter project", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req GetAppNamesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert platform strings to Platform types
	var platforms []namer.Platform
	if len(req.Platforms) > 0 {
		for _, platformStr := range req.Platforms {
			platform := namer.Platform(strings.ToLower(platformStr))
			platforms = append(platforms, platform)
		}
	}

	// Get app names
	appNamesResult, err := namer.GetAppNames(result.ProjectPath, platforms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appNamesResult)
}

// handleSetAppNames handles POST requests to set app names
func (api *NamerAPI) handleSetAppNames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate Flutter project
	result, err := flutter.ValidateProject(".")
	if err != nil || !result.IsValid {
		http.Error(w, "Not a valid Flutter project", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req SetAppNamesAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to namer request format
	namerRequest := &namer.SetAppNameRequest{
		ProjectPath: result.ProjectPath,
		Universal:   req.Universal,
		Platforms:   make(map[namer.Platform]string),
	}

	// Convert platform map
	for platformStr, appName := range req.Platforms {
		platform := namer.Platform(strings.ToLower(platformStr))
		namerRequest.Platforms[platform] = appName
	}

	// Set app names
	if err := namer.SetAppNames(namerRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated app names to return
	updatedResult, err := namer.GetAppNames(result.ProjectPath, nil)
	if err != nil {
		// Still return success even if we can't get updated names
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	// Return updated app names
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"app_names": updatedResult,
	})
}

// handleGetPlatforms handles GET requests to retrieve available platforms
func (api *NamerAPI) handleGetPlatforms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate Flutter project
	result, err := flutter.ValidateProject(".")
	if err != nil || !result.IsValid {
		http.Error(w, "Not a valid Flutter project", http.StatusBadRequest)
		return
	}

	// Get all platforms info
	allPlatforms := []PlatformInfo{
		{
			ID:          "android",
			Name:        "Android",
			Description: "Android mobile platform",
		},
		{
			ID:          "ios",
			Name:        "iOS",
			Description: "iOS mobile platform",
		},
		{
			ID:          "macos",
			Name:        "macOS",
			Description: "macOS desktop platform",
		},
		{
			ID:          "linux",
			Name:        "Linux",
			Description: "Linux desktop platform",
		},
		{
			ID:          "windows",
			Name:        "Windows",
			Description: "Windows desktop platform",
		},
		{
			ID:          "web",
			Name:        "Web",
			Description: "Web platform",
		},
	}

	// Check availability for each platform
	var availablePlatforms []PlatformInfo
	for _, platform := range allPlatforms {
		// Check if platform directory exists
		available := helpers.IsPlatformDirectoryAvailable(result.ProjectPath, platform.ID)
		platform.Available = available

		if available {
			availablePlatforms = append(availablePlatforms, platform)
		}
	}

	response := PlatformsResponse{
		Available: availablePlatforms,
		All:       allPlatforms,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetupNamerAPIRoutes sets up namer API routes
func SetupNamerAPIRoutes(project *flutter.ValidationResult) {
	namerAPI := NewNamerAPI()
	namerAPI.RegisterRoutes(http.DefaultServeMux)
}
