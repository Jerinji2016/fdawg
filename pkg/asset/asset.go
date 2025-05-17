package asset

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/flutter"
)

const (
	// AssetDirName is the name of the directory where assets are stored
	AssetDirName = "assets"
	// AssetBackupDirName is the name of the directory where assets are backed up during migration
	AssetBackupDirName = "assets.backup"
)

// AssetType represents the type of asset
type AssetType string

const (
	// ImageAsset represents image assets
	ImageAsset AssetType = "images"
	// AnimationAsset represents animation assets
	AnimationAsset AssetType = "animations"
	// AudioAsset represents audio assets
	AudioAsset AssetType = "audio"
	// VideoAsset represents video assets
	VideoAsset AssetType = "videos"
	// JSONAsset represents JSON assets
	JSONAsset AssetType = "json"
	// SVGAsset represents SVG assets
	SVGAsset AssetType = "svgs"
	// MiscAsset represents miscellaneous assets that don't fit in other categories
	MiscAsset AssetType = "misc"
)

// GetAssetDir returns the path to the asset directory for a Flutter project
func GetAssetDir(projectPath string) string {
	return filepath.Join(projectPath, AssetDirName)
}

// GetAssetBackupDir returns the path to the asset backup directory for a Flutter project
func GetAssetBackupDir(projectPath string) string {
	return filepath.Join(projectPath, AssetBackupDirName)
}

// EnsureAssetDirExists ensures that the main asset directory exists
func EnsureAssetDirExists(projectPath string) error {
	assetDir := GetAssetDir(projectPath)
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(assetDir, 0755); err != nil {
			return fmt.Errorf("failed to create asset directory: %v", err)
		}
	}
	return nil
}

// EnsureAssetTypeDirExists ensures that a specific asset type directory exists
func EnsureAssetTypeDirExists(projectPath string, assetType AssetType) error {
	// First ensure the main asset directory exists
	if err := EnsureAssetDirExists(projectPath); err != nil {
		return err
	}

	// Create the specific asset type directory
	assetTypeDir := filepath.Join(GetAssetDir(projectPath), string(assetType))
	if _, err := os.Stat(assetTypeDir); os.IsNotExist(err) {
		if err := os.MkdirAll(assetTypeDir, 0755); err != nil {
			return fmt.Errorf("failed to create asset type directory: %v", err)
		}
	}

	return nil
}

// DetermineAssetType determines the type of asset based on its file extension and content
func DetermineAssetType(filePath string) AssetType {
	ext := strings.ToLower(filepath.Ext(filePath))
	fileName := strings.ToLower(filepath.Base(filePath))

	// Image files
	if isInList(ext, []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".tiff", ".ico", ".heic", ".heif"}) {
		return ImageAsset
	}

	// Animation files
	if isInList(ext, []string{".flr", ".riv"}) ||
		(ext == ".json" && (strings.Contains(fileName, "animation") || strings.Contains(fileName, "lottie"))) {
		return AnimationAsset
	}

	// SVG files
	if ext == ".svg" {
		return SVGAsset
	}

	// Audio files
	if isInList(ext, []string{".mp3", ".wav", ".ogg", ".aac", ".flac", ".m4a", ".opus"}) {
		return AudioAsset
	}

	// Video files
	if isInList(ext, []string{".mp4", ".webm", ".avi", ".mov", ".mkv", ".flv", ".wmv", ".m4v"}) {
		return VideoAsset
	}

	// JSON files
	if ext == ".json" {
		return JSONAsset
	}

	// If no specific type is determined, categorize as miscellaneous
	return MiscAsset
}

// isInList checks if a string is in a list of strings
func isInList(str string, list []string) bool {
	return slices.Contains(list, str)
}

// AddAsset adds an asset to the project
func AddAsset(projectPath, assetPath string, assetType AssetType) error {
	// Check if the asset file exists
	if _, err := os.Stat(assetPath); os.IsNotExist(err) {
		return fmt.Errorf("asset file does not exist: %s", assetPath)
	}

	// If assetType is empty, determine it from the file extension
	if assetType == "" {
		assetType = DetermineAssetType(assetPath)
	}

	// Ensure the specific asset type directory exists
	if err := EnsureAssetTypeDirExists(projectPath, assetType); err != nil {
		return fmt.Errorf("failed to ensure asset type directory exists: %v", err)
	}

	// Get the destination directory
	destDir := filepath.Join(GetAssetDir(projectPath), string(assetType))

	// Copy the asset file to the destination directory
	destPath := filepath.Join(destDir, filepath.Base(assetPath))
	if err := copyFile(assetPath, destPath); err != nil {
		return fmt.Errorf("failed to copy asset file: %v", err)
	}

	// Update the pubspec.yaml file
	if err := updatePubspecWithAsset(projectPath, assetType); err != nil {
		return fmt.Errorf("failed to update pubspec.yaml: %v", err)
	}

	return nil
}

// RemoveAsset removes an asset from the project
func RemoveAsset(projectPath, assetName string, assetType AssetType) error {
	// Ensure the asset directory exists
	assetDir := GetAssetDir(projectPath)
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		return fmt.Errorf("asset directory does not exist")
	}

	// If assetType is empty, search for the asset in all directories
	if assetType == "" {
		assetTypes := []AssetType{
			ImageAsset,
			AnimationAsset,
			AudioAsset,
			VideoAsset,
			JSONAsset,
			SVGAsset,
			MiscAsset,
		}

		for _, at := range assetTypes {
			assetPath := filepath.Join(assetDir, string(at), assetName)
			if _, err := os.Stat(assetPath); err == nil {
				assetType = at
				break
			}
		}

		if assetType == "" {
			return fmt.Errorf("asset not found: %s", assetName)
		}
	}

	// Remove the asset file
	assetPath := filepath.Join(assetDir, string(assetType), assetName)
	if _, err := os.Stat(assetPath); os.IsNotExist(err) {
		return fmt.Errorf("asset file does not exist: %s", assetPath)
	}

	if err := os.Remove(assetPath); err != nil {
		return fmt.Errorf("failed to remove asset file: %v", err)
	}

	// Update the pubspec.yaml file
	if err := updatePubspecWithAsset(projectPath, assetType); err != nil {
		return fmt.Errorf("failed to update pubspec.yaml: %v", err)
	}

	return nil
}

// ListAssets lists all assets in the project
func ListAssets(projectPath string) (map[AssetType][]string, error) {
	// Ensure the asset directory exists
	assetDir := GetAssetDir(projectPath)
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("asset directory does not exist")
	}

	// List all assets by type
	assets := make(map[AssetType][]string)
	assetTypes := []AssetType{
		ImageAsset,
		AnimationAsset,
		AudioAsset,
		VideoAsset,
		JSONAsset,
		SVGAsset,
		MiscAsset,
	}

	for _, assetType := range assetTypes {
		assetTypeDir := filepath.Join(assetDir, string(assetType))
		if _, err := os.Stat(assetTypeDir); os.IsNotExist(err) {
			continue
		}

		files, err := os.ReadDir(assetTypeDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read asset directory: %v", err)
		}

		var assetFiles []string
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			assetFiles = append(assetFiles, file.Name())
		}

		assets[assetType] = assetFiles
	}

	return assets, nil
}

// MigrateAssets migrates assets from a flat structure to organized folders
func MigrateAssets(projectPath string) error {
	assetDir := GetAssetDir(projectPath)
	backupDir := GetAssetBackupDir(projectPath)

	// Check if the asset directory exists
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		return fmt.Errorf("asset directory does not exist")
	}

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Find all files in the asset directory (including in subdirectories)
	var filesToMigrate []string
	err := filepath.Walk(assetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == assetDir {
			return nil
		}

		// Skip directories that match our asset type directories
		relPath, err := filepath.Rel(assetDir, path)
		if err != nil {
			return err
		}

		// If it's a directory, check if it's one of our asset type directories
		if info.IsDir() {
			for _, assetType := range []AssetType{
				ImageAsset,
				AnimationAsset,
				AudioAsset,
				VideoAsset,
				JSONAsset,
				SVGAsset,
				MiscAsset,
			} {
				if relPath == string(assetType) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// It's a file, add it to the list
		filesToMigrate = append(filesToMigrate, path)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk asset directory: %v", err)
	}

	if len(filesToMigrate) == 0 {
		return fmt.Errorf("no files to migrate")
	}

	// Move all files to the backup directory
	for _, filePath := range filesToMigrate {
		// Get the relative path from the asset directory
		relPath, err := filepath.Rel(assetDir, filePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		// Create the backup path
		destPath := filepath.Join(backupDir, filepath.Base(filePath))

		// Copy the file to the backup directory
		if err := copyFile(filePath, destPath); err != nil {
			return fmt.Errorf("failed to backup file %s: %v", relPath, err)
		}
	}

	// Now add each file to the appropriate directory
	for _, filePath := range filesToMigrate {
		// Determine the asset type
		assetType := DetermineAssetType(filePath)

		// Ensure the asset type directory exists
		if err := EnsureAssetTypeDirExists(projectPath, assetType); err != nil {
			return fmt.Errorf("failed to create asset type directory: %v", err)
		}

		// Get the file name
		fileName := filepath.Base(filePath)

		// Copy the file to the appropriate directory
		destPath := filepath.Join(assetDir, string(assetType), fileName)
		if err := copyFile(filePath, destPath); err != nil {
			return fmt.Errorf("failed to copy file %s to %s: %v", fileName, assetType, err)
		}

		// Remove the original file
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to remove original file %s: %v", filePath, err)
		}
	}

	// Remove empty directories
	removeEmptyDirs(assetDir)

	// Update the pubspec.yaml file for each asset type
	assetTypes := []AssetType{
		ImageAsset,
		AnimationAsset,
		AudioAsset,
		VideoAsset,
		JSONAsset,
		SVGAsset,
		MiscAsset,
	}

	for _, assetType := range assetTypes {
		assetTypeDir := filepath.Join(assetDir, string(assetType))
		if _, err := os.Stat(assetTypeDir); os.IsNotExist(err) {
			continue
		}

		if err := updatePubspecWithAsset(projectPath, assetType); err != nil {
			return fmt.Errorf("failed to update pubspec.yaml for %s: %v", assetType, err)
		}
	}

	// Generate the Dart asset file
	if err := GenerateDartAssetFile(projectPath); err != nil {
		return fmt.Errorf("failed to generate Dart asset file: %v", err)
	}

	return nil
}

// removeEmptyDirs removes empty directories recursively
func removeEmptyDirs(dir string) error {
	// Skip asset type directories
	assetTypes := []AssetType{
		ImageAsset,
		AnimationAsset,
		AudioAsset,
		VideoAsset,
		JSONAsset,
		SVGAsset,
		MiscAsset,
	}

	for _, assetType := range assetTypes {
		if filepath.Base(dir) == string(assetType) {
			return nil
		}
	}

	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Process each entry
	for _, entry := range entries {
		if entry.IsDir() {
			// Recursively process subdirectories
			subDir := filepath.Join(dir, entry.Name())
			if err := removeEmptyDirs(subDir); err != nil {
				return err
			}
		}
	}

	// Check if directory is empty now
	entries, err = os.ReadDir(dir)
	if err != nil {
		return err
	}

	// If directory is empty, remove it
	if len(entries) == 0 {
		if err := os.Remove(dir); err != nil {
			return err
		}
	}

	return nil
}

// updatePubspecWithAsset is implemented in pubspec_updater.go

// GenerateDartAssetFile generates a Dart asset file with all assets
func GenerateDartAssetFile(projectPath string) error {
	// List all assets
	assets, err := ListAssets(projectPath)
	if err != nil {
		return fmt.Errorf("failed to list assets: %v", err)
	}

	// Create the Dart file content
	var content strings.Builder

	// Add file header
	content.WriteString(`// GENERATED CODE - DO NOT MODIFY BY HAND
// Generated by fdawg

/// Asset management class
///
/// This class provides access to assets defined in the assets directory.
/// It is automatically generated by fdawg and should not be modified manually.
class Asset {
  // Private constructor to prevent instantiation
  Asset._();

  /// Images assets
  static final Images images = Images._();

  /// Animations assets
  static final Animations animations = Animations._();

  /// Audio assets
  static final Audio audio = Audio._();

  /// Video assets
  static final Videos videos = Videos._();

  /// JSON assets
  static final Json json = Json._();

  /// SVG assets
  static final Svgs svgs = Svgs._();

  /// Miscellaneous assets
  static final Misc misc = Misc._();
}

`)

	// Add asset classes
	addAssetClass(&content, "Images", ImageAsset, assets[ImageAsset])
	addAssetClass(&content, "Animations", AnimationAsset, assets[AnimationAsset])
	addAssetClass(&content, "Audio", AudioAsset, assets[AudioAsset])
	addAssetClass(&content, "Videos", VideoAsset, assets[VideoAsset])
	addAssetClass(&content, "Json", JSONAsset, assets[JSONAsset])
	addAssetClass(&content, "Svgs", SVGAsset, assets[SVGAsset])
	addAssetClass(&content, "Misc", MiscAsset, assets[MiscAsset])

	// Write the file
	dartFilePath := filepath.Join(projectPath, "lib", "config")

	// Ensure the directory exists
	if err := os.MkdirAll(dartFilePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	dartFilePath = filepath.Join(dartFilePath, "asset.dart")

	if err := os.WriteFile(dartFilePath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write Dart file: %v", err)
	}

	return nil
}

// addAssetClass adds an asset class to the content builder
func addAssetClass(content *strings.Builder, className string, assetType AssetType, assetFiles []string) {
	content.WriteString(fmt.Sprintf(`/// %s assets
class %s {
  // Private constructor to prevent instantiation
  %s._();

`, className, className, className))

	// Add constants for each asset
	for _, assetFile := range assetFiles {
		// Create a valid Dart variable name
		varName := createDartVariableName(assetFile)

		// Add the constant
		content.WriteString(fmt.Sprintf(`  /// %s asset
  static const String %s = 'assets/%s/%s';

`, assetFile, varName, assetType, assetFile))
	}

	content.WriteString("}\n\n")
}

// createDartVariableName creates a valid Dart variable name from a file name
func createDartVariableName(fileName string) string {
	// Remove the file extension
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// Use the common helper function to format the name
	return flutter.FormatDartVariableName(name)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
