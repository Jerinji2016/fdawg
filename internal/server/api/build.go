package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	mux.HandleFunc("/api/build/config", api.handleGetConfig)
	mux.HandleFunc("/api/build/config/update", api.handleUpdateConfig)
	mux.HandleFunc("/api/build/stream", api.handleBuildStream)
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
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Platform  string    `json:"platform"`
	Type      string    `json:"type"`
	Size      string    `json:"size"`
	Date      string    `json:"date"`
	Timestamp time.Time `json:"timestamp"`
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

	artifacts, err := api.scanBuildArtifacts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to scan artifacts: %v", err), http.StatusInternalServerError)
		return
	}

	response := BuildArtifactsResponse{
		Artifacts: artifacts,
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

	relativePath := r.URL.Query().Get("path")
	if relativePath == "" {
		http.Error(w, "Path parameter required", http.StatusBadRequest)
		return
	}

	// Get build config to determine output directory
	configPath := ".fdawg/build.yaml"
	if !api.buildConfigExists(configPath) {
		http.Error(w, "Build configuration not found", http.StatusNotFound)
		return
	}

	buildConfig, err := build.LoadBuildConfig(api.project.ProjectPath, configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load build config: %v", err), http.StatusInternalServerError)
		return
	}

	outputDir := buildConfig.Artifacts.BaseOutputDir
	if outputDir == "" {
		outputDir = "output"
	}

	// Construct full file path
	fullPath := filepath.Join(api.project.ProjectPath, outputDir, relativePath)

	// Security check: ensure the file is within the output directory
	outputPath := filepath.Join(api.project.ProjectPath, outputDir)
	if !strings.HasPrefix(fullPath, outputPath) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Check if file exists
	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("File access error: %v", err), http.StatusInternalServerError)
		return
	}

	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set headers for download
	filename := filepath.Base(fullPath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Stream the file to the response
	http.ServeContent(w, r, filename, fileInfo.ModTime(), file)
}

// handleGetConfig handles GET requests to get build configuration for editing
func (api *BuildAPI) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	configPath := ".fdawg/build.yaml"
	if !api.buildConfigExists(configPath) {
		http.Error(w, "Build configuration not found", http.StatusNotFound)
		return
	}

	buildConfig, err := build.LoadBuildConfig(api.project.ProjectPath, configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load build config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildConfig)
}

// handleUpdateConfig handles POST requests to update build configuration
func (api *BuildAPI) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config build.BuildConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	configPath := ".fdawg/build.yaml"

	// Save the updated configuration
	if err := build.SaveBuildConfig(api.project.ProjectPath, configPath, &config); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save configuration: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleBuildStream handles Server-Sent Events for real-time build streaming
func (api *BuildAPI) handleBuildStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Create a flusher to send data immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send initial connection message
	fmt.Fprintf(w, "data: %s\n\n", `{"type":"log","message":"Build stream connected","level":"info"}`)
	flusher.Flush()

	// Keep connection alive and send periodic heartbeats
	// For now, this is a simple implementation
	// In a real implementation, you would:
	// 1. Subscribe to build events from the build manager
	// 2. Stream real-time build logs and progress
	// 3. Handle client disconnection properly

	// Simple heartbeat for demonstration
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Listen for client disconnect
	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			// Client disconnected
			return
		case <-ticker.C:
			// Send heartbeat
			fmt.Fprintf(w, "data: %s\n\n", `{"type":"log","message":"Connection alive","level":"debug"}`)
			flusher.Flush()
		}
	}
}

// Helper methods

func (api *BuildAPI) buildConfigExists(configPath string) bool {
	fullPath := filepath.Join(api.project.ProjectPath, configPath)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (api *BuildAPI) getLastBuildInfo() string {
	// Check for recent artifacts to determine last build
	artifacts, err := api.scanBuildArtifacts()
	if err != nil || len(artifacts) == 0 {
		return "No recent builds"
	}

	// Find the most recent artifact
	var mostRecent *BuildArtifactInfo
	for i := range artifacts {
		if mostRecent == nil {
			mostRecent = &artifacts[i]
			continue
		}
		// Simple comparison by date string (could be improved with actual time parsing)
		if artifacts[i].Date > mostRecent.Date {
			mostRecent = &artifacts[i]
		}
	}

	if mostRecent != nil {
		return fmt.Sprintf("Last build: %s (%s)", mostRecent.Date, mostRecent.Platform)
	}

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

func (api *BuildAPI) scanBuildArtifacts() ([]BuildArtifactInfo, error) {
	var artifacts []BuildArtifactInfo

	// Check if build config exists to get output directory
	configPath := ".fdawg/build.yaml"
	if !api.buildConfigExists(configPath) {
		return artifacts, nil // Return empty if no config
	}

	buildConfig, err := build.LoadBuildConfig(api.project.ProjectPath, configPath)
	if err != nil {
		return artifacts, err
	}

	outputDir := buildConfig.Artifacts.BaseOutputDir
	if outputDir == "" {
		outputDir = "output" // Default output directory
	}

	outputPath := filepath.Join(api.project.ProjectPath, outputDir)

	// Check if output directory exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return artifacts, nil // Return empty if no output directory
	}

	// Walk through the output directory
	err = filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include build artifacts (common extensions)
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if !isArtifactFile(ext) {
			return nil
		}

		// Determine platform from path
		platform := determinePlatformFromPath(path, outputPath)

		// Get relative path from output directory
		relPath, _ := filepath.Rel(outputPath, path)

		// Format file size
		size := formatFileSize(info.Size())

		// Format date
		date := info.ModTime().Format("Jan 2, 2006 15:04")

		artifact := BuildArtifactInfo{
			Name:      info.Name(),
			Path:      relPath,
			Platform:  platform,
			Type:      getArtifactType(ext),
			Size:      size,
			Date:      date,
			Timestamp: info.ModTime(),
		}

		artifacts = append(artifacts, artifact)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort artifacts by date (newest first)
	sortArtifactsByDate(artifacts)

	return artifacts, nil
}

// sortArtifactsByDate sorts artifacts by creation time (newest first)
func sortArtifactsByDate(artifacts []BuildArtifactInfo) {
	// Use sort.Slice to sort by timestamp in descending order (newest first)
	for i := 0; i < len(artifacts); i++ {
		for j := i + 1; j < len(artifacts); j++ {
			if artifacts[i].Timestamp.Before(artifacts[j].Timestamp) {
				artifacts[i], artifacts[j] = artifacts[j], artifacts[i]
			}
		}
	}
}

func isArtifactFile(ext string) bool {
	artifactExtensions := []string{
		".apk", ".aab", ".ipa", ".app", ".dmg", ".exe", ".msix",
		".deb", ".rpm", ".tar.gz", ".zip", ".tar.xz",
	}

	for _, validExt := range artifactExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func determinePlatformFromPath(filePath, outputPath string) string {
	relPath, _ := filepath.Rel(outputPath, filePath)
	pathParts := strings.Split(relPath, string(filepath.Separator))

	// Look for platform indicators in path
	for _, part := range pathParts {
		part = strings.ToLower(part)
		switch {
		case strings.Contains(part, "android"):
			return "android"
		case strings.Contains(part, "ios"):
			return "ios"
		case strings.Contains(part, "macos"):
			return "macos"
		case strings.Contains(part, "linux"):
			return "linux"
		case strings.Contains(part, "windows"):
			return "windows"
		case strings.Contains(part, "web"):
			return "web"
		}
	}

	// Fallback: determine by file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".apk", ".aab":
		return "android"
	case ".ipa":
		return "ios"
	case ".app", ".dmg":
		return "macos"
	case ".exe", ".msix":
		return "windows"
	case ".deb", ".rpm", ".tar.gz", ".tar.xz":
		return "linux"
	default:
		return "unknown"
	}
}

func getArtifactType(ext string) string {
	switch ext {
	case ".apk":
		return "Android APK"
	case ".aab":
		return "Android Bundle"
	case ".ipa":
		return "iOS App"
	case ".app":
		return "macOS App"
	case ".dmg":
		return "macOS Installer"
	case ".exe":
		return "Windows Executable"
	case ".msix":
		return "Windows Package"
	case ".deb":
		return "Debian Package"
	case ".rpm":
		return "RPM Package"
	case ".tar.gz", ".tar.xz":
		return "Linux Archive"
	case ".zip":
		return "Archive"
	default:
		return "Build Artifact"
	}
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
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
