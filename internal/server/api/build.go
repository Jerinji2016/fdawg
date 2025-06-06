package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jerinji2016/fdawg/internal/server/helpers"
	"github.com/Jerinji2016/fdawg/pkg/build"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
)

// BuildAPI handles build-related API endpoints
type BuildAPI struct {
	project *flutter.ValidationResult
}

// NewBuildAPI creates a new BuildAPI instance
func NewBuildAPI(project *flutter.ValidationResult) *BuildAPI {
	return &BuildAPI{
		project: project,
	}
}

// RegisterRoutes registers build API routes
func (api *BuildAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/build/status", api.handleGetStatus)
	mux.HandleFunc("/api/build/setup", api.handleSetup)
	mux.HandleFunc("/api/build/run", api.handleRun)
	mux.HandleFunc("/api/build/stop", api.handleStop)
	mux.HandleFunc("/api/build/reset", api.handleReset)
	mux.HandleFunc("/api/build/platforms", api.handleGetPlatforms)
	mux.HandleFunc("/api/build/artifacts", api.handleGetArtifacts)
	mux.HandleFunc("/api/build/artifacts/download", api.handleDownloadArtifact)
}

// Request/Response types for build API

type BuildStatusResponse struct {
	ConfigExists bool               `json:"config_exists"`
	Config       *build.BuildConfig `json:"config,omitempty"`
	LastBuild    string             `json:"last_build,omitempty"`
	Error        string             `json:"error,omitempty"`
}

type BuildSetupRequest struct {
	Default bool `json:"default"`
	Force   bool `json:"force"`
}

type BuildRunRequest struct {
	Platforms       []string `json:"platforms"`
	Environment     string   `json:"environment,omitempty"`
	SkipPreBuild    bool     `json:"skip_pre_build"`
	ContinueOnError bool     `json:"continue_on_error"`
	DryRun          bool     `json:"dry_run"`
	Parallel        bool     `json:"parallel"`
}

type BuildArtifactInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Platform string `json:"platform"`
	Type     string `json:"type"`
	Size     string `json:"size"`
	Date     string `json:"date"`
}

type BuildArtifactsResponse struct {
	Artifacts []BuildArtifactInfo `json:"artifacts"`
}

// handleGetStatus handles GET requests to get build status
func (api *BuildAPI) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	configPath := ".fdawg/build.yaml"
	configExists := api.buildConfigExists(configPath)

	response := BuildStatusResponse{
		ConfigExists: configExists,
	}

	if configExists {
		// Load configuration
		buildConfig, err := build.LoadBuildConfig(api.project.ProjectPath, configPath)
		if err != nil {
			response.Error = fmt.Sprintf("Failed to load config: %v", err)
		} else {
			response.Config = buildConfig
		}

		// Get last build info (simplified for now)
		response.LastBuild = api.getLastBuildInfo()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSetup handles POST requests to setup build configuration
func (api *BuildAPI) handleSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BuildSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	configPath := ".fdawg/build.yaml"

	// Check if configuration already exists
	if api.buildConfigExists(configPath) && !req.Force {
		http.Error(w, "Build configuration already exists", http.StatusConflict)
		return
	}

	var config *build.BuildConfig

	if req.Default {
		// Create default configuration
		config = api.createDefaultConfigForProject()
	} else {
		// For now, we'll also use default config for wizard
		// In the future, this could be an interactive process
		config = api.createDefaultConfigForProject()
	}

	// Save configuration
	if err := build.SaveBuildConfig(api.project.ProjectPath, configPath, config); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save configuration: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleRun handles POST requests to run builds
func (api *BuildAPI) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BuildRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	configPath := ".fdawg/build.yaml"
	if !api.buildConfigExists(configPath) {
		http.Error(w, "Build configuration not found. Please run setup first.", http.StatusBadRequest)
		return
	}

	// Convert platform strings to build.Platform
	var platforms []build.Platform
	for _, platformStr := range req.Platforms {
		platform := build.Platform(strings.ToLower(platformStr))
		platforms = append(platforms, platform)
	}

	if len(platforms) == 0 {
		http.Error(w, "No platforms specified", http.StatusBadRequest)
		return
	}

	// Load build configuration
	buildConfig, err := build.LoadBuildConfig(api.project.ProjectPath, configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load build config: %v", err), http.StatusInternalServerError)
		return
	}

	// Validate environment if specified
	if req.Environment != "" {
		if err := api.validateEnvironment(req.Environment); err != nil {
			http.Error(w, fmt.Sprintf("Environment validation failed: %v", err), http.StatusBadRequest)
			return
		}
	}

	// Create build manager
	buildManager, err := build.NewBuildManager(api.project.ProjectPath, buildConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create build manager: %v", err), http.StatusInternalServerError)
		return
	}

	// Build options
	options := build.BuildOptions{
		SkipPreBuild:    req.SkipPreBuild,
		ContinueOnError: req.ContinueOnError,
		DryRun:          req.DryRun,
		Parallel:        req.Parallel,
		Environment:     req.Environment,
	}

	if options.DryRun {
		// Show build plan
		if err := buildManager.ShowBuildPlan(platforms, options); err != nil {
			http.Error(w, fmt.Sprintf("Failed to show build plan: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"dry_run": true,
			"message": "Build plan shown in logs",
		})
		return
	}

	// Execute build
	result, err := buildManager.ExecuteBuild(platforms, options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Build failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleStop handles POST requests to stop builds
func (api *BuildAPI) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, just return success
	// In the future, this could actually stop running builds
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleReset handles POST requests to reset build configuration
func (api *BuildAPI) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	configPath := filepath.Join(api.project.ProjectPath, ".fdawg", "build.yaml")

	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("Failed to remove configuration: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleGetPlatforms handles GET requests to get platform information
func (api *BuildAPI) handleGetPlatforms(w http.ResponseWriter, r *http.Request) {
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

// handleGetArtifacts handles GET requests to get build artifacts
func (api *BuildAPI) handleGetArtifacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, return empty artifacts
	// In the future, this would scan the build output directory
	response := BuildArtifactsResponse{
		Artifacts: []BuildArtifactInfo{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleDownloadArtifact handles GET requests to download build artifacts
func (api *BuildAPI) handleDownloadArtifact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path parameter required", http.StatusBadRequest)
		return
	}

	// For now, just return an error
	// In the future, this would serve the actual file
	http.Error(w, "Download not yet implemented", http.StatusNotImplemented)
}

// Helper methods

func (api *BuildAPI) buildConfigExists(configPath string) bool {
	fullPath := filepath.Join(api.project.ProjectPath, configPath)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (api *BuildAPI) getLastBuildInfo() string {
	// For now, return a placeholder
	// In the future, this would check build logs or artifacts
	return "No recent builds"
}

func (api *BuildAPI) createDefaultConfigForProject() *build.BuildConfig {
	// Create a default configuration similar to the CLI version
	config := build.DefaultBuildConfig()

	// Customize based on available platforms
	availablePlatforms := []string{}
	platforms := []string{"android", "ios", "web", "macos", "linux", "windows"}

	for _, platform := range platforms {
		if helpers.IsPlatformDirectoryAvailable(api.project.ProjectPath, platform) {
			availablePlatforms = append(availablePlatforms, platform)
		}
	}

	// Enable only available platforms
	config.Platforms.Android.Enabled = contains(availablePlatforms, "android")
	config.Platforms.IOS.Enabled = contains(availablePlatforms, "ios")
	config.Platforms.Web.Enabled = contains(availablePlatforms, "web")
	config.Platforms.MacOS.Enabled = contains(availablePlatforms, "macos")
	config.Platforms.Linux.Enabled = contains(availablePlatforms, "linux")
	config.Platforms.Windows.Enabled = contains(availablePlatforms, "windows")

	return config
}

func (api *BuildAPI) validateEnvironment(envName string) error {
	// For now, just return nil
	// In the future, this would validate the environment exists
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SetupBuildAPIRoutes sets up build API routes
func SetupBuildAPIRoutes(project *flutter.ValidationResult) {
	buildAPI := NewBuildAPI(project)
	buildAPI.RegisterRoutes(http.DefaultServeMux)
}
