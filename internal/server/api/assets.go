package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Jerinji2016/fdawg/pkg/asset"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
)

// SetupAssetAPIRoutes sets up the API routes for asset management
func SetupAssetAPIRoutes(project *flutter.ValidationResult) {
	// Upload asset API route
	http.HandleFunc("/api/assets/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form data (max 32MB)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to parse form data: %v", err),
			})
			return
		}

		// Get asset type
		assetType := asset.AssetType(r.FormValue("asset_type"))
		if assetType == "" {
			// Default to auto-detection
			assetType = ""
		}

		// Get uploaded files
		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "No files uploaded",
			})
			return
		}

		// Process each file
		results := make([]map[string]interface{}, 0)
		for _, fileHeader := range files {
			// Open the uploaded file
			file, err := fileHeader.Open()
			if err != nil {
				results = append(results, map[string]interface{}{
					"filename": fileHeader.Filename,
					"success":  false,
					"error":    fmt.Sprintf("Failed to open file: %v", err),
				})
				continue
			}
			defer file.Close()

			// Create a temporary directory if it doesn't exist
			tempDir := filepath.Join(os.TempDir(), "fdawg-assets")
			if err := os.MkdirAll(tempDir, 0755); err != nil {
				results = append(results, map[string]interface{}{
					"filename": fileHeader.Filename,
					"success":  false,
					"error":    fmt.Sprintf("Failed to create temporary directory: %v", err),
				})
				continue
			}

			// Create a temporary file with the original filename
			tempFilePath := filepath.Join(tempDir, fileHeader.Filename)
			tempFile, err := os.Create(tempFilePath)
			if err != nil {
				results = append(results, map[string]interface{}{
					"filename": fileHeader.Filename,
					"success":  false,
					"error":    fmt.Sprintf("Failed to create temporary file: %v", err),
				})
				continue
			}
			defer os.Remove(tempFilePath)
			defer tempFile.Close()

			// Copy the uploaded file to the temporary file
			_, err = io.Copy(tempFile, file)
			if err != nil {
				results = append(results, map[string]interface{}{
					"filename": fileHeader.Filename,
					"success":  false,
					"error":    fmt.Sprintf("Failed to save file: %v", err),
				})
				continue
			}

			// Determine asset type if not specified
			currentAssetType := assetType
			if currentAssetType == "" {
				currentAssetType = asset.DetermineAssetType(fileHeader.Filename)
			}

			// Make sure the file is flushed to disk
			tempFile.Sync()
			tempFile.Close()

			// Add the asset
			err = asset.AddAsset(project.ProjectPath, tempFilePath, currentAssetType)
			if err != nil {
				results = append(results, map[string]interface{}{
					"filename": fileHeader.Filename,
					"success":  false,
					"error":    fmt.Sprintf("Failed to add asset: %v", err),
				})
				continue
			}

			// Add success result
			results = append(results, map[string]interface{}{
				"filename":   fileHeader.Filename,
				"success":    true,
				"asset_type": string(currentAssetType),
			})
		}

		// Generate the Dart asset file
		err = asset.GenerateDartAssetFile(project.ProjectPath)
		if err != nil {
			// Just log the error, don't fail the whole operation
			fmt.Printf("Warning: Failed to generate Dart asset file: %v\n", err)
		}

		// Return results as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"results": results,
		})
	})

	// Delete asset API route
	http.HandleFunc("/api/assets/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Set content type for all responses
		w.Header().Set("Content-Type", "application/json")

		// Parse multipart form data
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			// Try regular form parsing if multipart fails
			err = r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   fmt.Sprintf("Failed to parse form data: %v", err),
				})
				return
			}
		}

		// Log the form data for debugging
		fmt.Printf("Delete asset form data: %+v\n", r.Form)

		// Get asset name and type
		assetName := r.FormValue("asset_name")
		if assetName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Asset name is required",
			})
			return
		}

		assetType := asset.AssetType(r.FormValue("asset_type"))

		// Log the values for debugging
		fmt.Printf("Deleting asset: name=%s, type=%s\n", assetName, assetType)

		// Remove the asset
		err = asset.RemoveAsset(project.ProjectPath, assetName, assetType)
		if err != nil {
			// Set content type to JSON before sending error
			w.Header().Set("Content-Type", "application/json")
			// Return error as JSON
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to remove asset: %v", err),
			})
			return
		}

		// Generate the Dart asset file
		err = asset.GenerateDartAssetFile(project.ProjectPath)
		if err != nil {
			// Just log the error, don't fail the whole operation
			fmt.Printf("Warning: Failed to generate Dart asset file: %v\n", err)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	})

	// Download asset API route
	http.HandleFunc("/api/assets/download", func(w http.ResponseWriter, r *http.Request) {
		// Get asset name and type
		assetName := r.URL.Query().Get("asset_name")
		if assetName == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Asset name is required",
			})
			return
		}

		assetType := asset.AssetType(r.URL.Query().Get("asset_type"))
		if assetType == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Asset type is required",
			})
			return
		}

		// Get asset path
		assetDir := asset.GetAssetDir(project.ProjectPath)
		assetPath := filepath.Join(assetDir, string(assetType), assetName)

		// Check if asset exists
		if _, err := os.Stat(assetPath); os.IsNotExist(err) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Asset not found",
			})
			return
		}

		// Set headers for file download
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", assetName))
		w.Header().Set("Content-Type", "application/octet-stream")

		// Serve the file
		http.ServeFile(w, r, assetPath)
	})

	// List assets API route
	http.HandleFunc("/api/assets/list", func(w http.ResponseWriter, r *http.Request) {
		// List assets
		assets, err := asset.ListAssets(project.ProjectPath)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to list assets: %v", err),
			})
			return
		}

		// Return assets as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"assets":  assets,
		})
	})

	// Migrate assets API route
	http.HandleFunc("/api/assets/migrate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Migrate assets
		err := asset.MigrateAssets(project.ProjectPath)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to migrate assets: %v", err),
			})
			return
		}

		// Generate the Dart asset file
		err = asset.GenerateDartAssetFile(project.ProjectPath)
		if err != nil {
			// Just log the error, don't fail the whole operation
			fmt.Printf("Warning: Failed to generate Dart asset file: %v\n", err)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	})
}
