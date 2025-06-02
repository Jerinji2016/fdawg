package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Jerinji2016/fdawg/internal/server/helpers"
	"github.com/Jerinji2016/fdawg/pkg/bundler"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
)

// BundlerAPI handles bundler-related API endpoints
type BundlerAPI struct {
	project *flutter.ValidationResult
}

// NewBundlerAPI creates a new BundlerAPI instance
func NewBundlerAPI(project *flutter.ValidationResult) *BundlerAPI {
	return &BundlerAPI{
		project: project,
	}
}

// RegisterRoutes registers bundler API routes
func (api *BundlerAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/bundler/get", api.handleGetBundleIDs)
	mux.HandleFunc("/api/bundler/set", api.handleSetBundleIDs)
	mux.HandleFunc("/api/bundler/platforms", api.handleGetPlatforms)
	mux.HandleFunc("/api/bundler/validate", api.handleValidateBundleID)
}

// Request/Response types for bundler API

type GetBundleIDsAPIRequest struct {
	Platforms []string `json:"platforms,omitempty"`
}

type SetBundleIDsAPIRequest struct {
	Universal string            `json:"universal,omitempty"`
	Platforms map[string]string `json:"platforms,omitempty"`
}

type ValidateBundleIDAPIRequest struct {
	BundleID string `json:"bundle_id"`
}

type ValidateBundleIDAPIResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// handleGetBundleIDs handles POST requests to get bundle IDs
func (api *BundlerAPI) handleGetBundleIDs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req GetBundleIDsAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert platform strings to bundler.Platform
	var platforms []bundler.Platform
	for _, platformStr := range req.Platforms {
		platform := bundler.Platform(strings.ToLower(platformStr))
		platforms = append(platforms, platform)
	}

	// Get bundle IDs
	bundleIDsResult, err := bundler.GetBundleIDs(api.project.ProjectPath, platforms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bundleIDsResult)
}

// handleSetBundleIDs handles POST requests to set bundle IDs
func (api *BundlerAPI) handleSetBundleIDs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req SetBundleIDsAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to bundler request format
	bundlerRequest := &bundler.BundleIDRequest{
		ProjectPath: api.project.ProjectPath,
		Universal:   req.Universal,
		Platforms:   make(map[bundler.Platform]string),
	}

	// Convert platform map
	for platformStr, bundleID := range req.Platforms {
		platform := bundler.Platform(strings.ToLower(platformStr))
		bundlerRequest.Platforms[platform] = bundleID
	}

	// Set bundle IDs
	if err := bundler.SetBundleIDs(bundlerRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated bundle IDs to return
	updatedResult, err := bundler.GetBundleIDs(api.project.ProjectPath, nil)
	if err != nil {
		// Still return success even if we can't get updated bundle IDs
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	// Return updated bundle IDs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedResult)
}

// handleValidateBundleID handles POST requests to validate bundle ID format
func (api *BundlerAPI) handleValidateBundleID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req ValidateBundleIDAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate bundle ID (we'll need to create a validation function)
	valid, errorMsg := validateBundleIDFormat(req.BundleID)

	response := ValidateBundleIDAPIResponse{
		Valid: valid,
		Error: errorMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetPlatforms handles GET requests to get platform information
func (api *BundlerAPI) handleGetPlatforms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define all platforms with their information
	allPlatforms := []PlatformInfo{
		{ID: "android", Name: "Android", Description: "Android mobile platform"},
		{ID: "ios", Name: "iOS", Description: "iOS mobile platform"},
		{ID: "macos", Name: "macOS", Description: "macOS desktop platform"},
		{ID: "linux", Name: "Linux", Description: "Linux desktop platform"},
		{ID: "windows", Name: "Windows", Description: "Windows desktop platform"},
		{ID: "web", Name: "Web", Description: "Web platform"},
	}

	// Check availability for each platform
	var availablePlatforms []PlatformInfo
	for i, platform := range allPlatforms {
		// Check if platform directory exists
		available := helpers.IsPlatformDirectoryAvailable(api.project.ProjectPath, platform.ID)
		allPlatforms[i].Available = available

		if available {
			availablePlatforms = append(availablePlatforms, allPlatforms[i])
		}
	}

	response := PlatformsResponse{
		Available: availablePlatforms,
		All:       allPlatforms,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// validateBundleIDFormat validates bundle ID format and naming conventions
func validateBundleIDFormat(bundleID string) (bool, string) {
	if bundleID == "" {
		return false, "Bundle ID cannot be empty"
	}

	// Check for reverse domain notation (at least one dot)
	if !strings.Contains(bundleID, ".") {
		return false, "Bundle ID should follow reverse domain notation (e.g., com.company.app)"
	}

	// Check for valid characters (alphanumeric, dots, hyphens, underscores)
	for _, char := range bundleID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '.' || char == '-' || char == '_') {
			return false, "Bundle ID contains invalid characters. Only alphanumeric characters, dots, hyphens, and underscores are allowed"
		}
	}

	// Check that it doesn't start or end with a dot
	if strings.HasPrefix(bundleID, ".") || strings.HasSuffix(bundleID, ".") {
		return false, "Bundle ID cannot start or end with a dot"
	}

	// Check for consecutive dots
	if strings.Contains(bundleID, "..") {
		return false, "Bundle ID cannot contain consecutive dots"
	}

	// Check minimum length
	if len(bundleID) < 3 {
		return false, "Bundle ID is too short (minimum 3 characters)"
	}

	// Check maximum length (reasonable limit)
	if len(bundleID) > 255 {
		return false, "Bundle ID is too long (maximum 255 characters)"
	}

	// Split by dots and validate each segment
	segments := strings.Split(bundleID, ".")
	for _, segment := range segments {
		if segment == "" {
			return false, "Bundle ID segment is empty"
		}

		// Each segment should not start with a number (Java package naming convention)
		if len(segment) > 0 && segment[0] >= '0' && segment[0] <= '9' {
			return false, "Bundle ID segment '" + segment + "' cannot start with a number"
		}
	}

	return true, ""
}

// SetupBundlerAPIRoutes sets up bundler API routes
func SetupBundlerAPIRoutes(project *flutter.ValidationResult) {
	bundlerAPI := NewBundlerAPI(project)
	bundlerAPI.RegisterRoutes(http.DefaultServeMux)
}
