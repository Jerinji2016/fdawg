package asset

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"gopkg.in/yaml.v3"
)

const (
	// AssetDirName is the name of the directory where assets are stored
	AssetDirName = "assets"
)

// AssetType represents the type of asset
type AssetType string

const (
	// ImageAsset represents image assets
	ImageAsset AssetType = "images"
	// FontAsset represents font assets
	FontAsset AssetType = "fonts"
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
	// TranslationAsset represents translation assets
	TranslationAsset AssetType = "translations"
)

// GetAssetDir returns the path to the asset directory for a Flutter project
func GetAssetDir(projectPath string) string {
	return filepath.Join(projectPath, AssetDirName)
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

// DetermineAssetType determines the type of asset based on its file extension
func DetermineAssetType(filePath string) AssetType {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp":
		return ImageAsset
	case ".ttf", ".otf":
		return FontAsset
	case ".json":
		// Check if it's a Lottie animation (could be more sophisticated)
		if strings.Contains(strings.ToLower(filePath), "animation") {
			return AnimationAsset
		}
		return JSONAsset
	case ".svg":
		return SVGAsset
	case ".mp3", ".wav", ".ogg", ".aac":
		return AudioAsset
	case ".mp4", ".webm", ".avi", ".mov":
		return VideoAsset
	case ".arb":
		return TranslationAsset
	default:
		// Default to image assets
		return ImageAsset
	}
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
			FontAsset,
			AnimationAsset,
			AudioAsset,
			VideoAsset,
			JSONAsset,
			SVGAsset,
			TranslationAsset,
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
		FontAsset,
		AnimationAsset,
		AudioAsset,
		VideoAsset,
		JSONAsset,
		SVGAsset,
		TranslationAsset,
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

// updatePubspecWithAsset updates the pubspec.yaml file with the asset entry
func updatePubspecWithAsset(projectPath string, assetType AssetType) error {
	// Read the pubspec.yaml file
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	pubspecData, err := os.ReadFile(pubspecPath)
	if err != nil {
		return fmt.Errorf("failed to read pubspec.yaml: %v", err)
	}

	// Create the asset path
	assetPath := fmt.Sprintf("assets/%s/", assetType)

	// Parse the pubspec.yaml file to check if the asset path already exists
	var pubspec map[string]interface{}
	if err := yaml.Unmarshal(pubspecData, &pubspec); err != nil {
		return fmt.Errorf("failed to parse pubspec.yaml: %v", err)
	}

	// Check if the flutter section exists
	flutter, ok := pubspec["flutter"].(map[string]interface{})
	if !ok {
		// If flutter section doesn't exist, we need to add it with the assets
		return updatePubspecWithNewFlutterSection(pubspecPath, string(pubspecData), assetPath)
	}

	// Check if the assets section exists
	assets, ok := flutter["assets"].([]interface{})
	if !ok {
		// If assets section doesn't exist, we need to add it
		return updatePubspecWithNewAssetsSection(pubspecPath, string(pubspecData), assetPath)
	}

	// Check if the asset path already exists
	for _, asset := range assets {
		if assetStr, ok := asset.(string); ok && assetStr == assetPath {
			// Asset path already exists, no need to update
			return nil
		}
	}

	// Asset path doesn't exist, we need to add it to the existing assets section
	return updatePubspecWithNewAssetPath(pubspecPath, string(pubspecData), assetPath)
}

// updatePubspecWithNewFlutterSection adds a new flutter section with assets to pubspec.yaml
func updatePubspecWithNewFlutterSection(pubspecPath, pubspecContent, assetPath string) error {
	// Find the end of the file or the line before another top-level section
	lines := strings.Split(pubspecContent, "\n")
	insertIndex := len(lines)

	// Look for a good place to insert the flutter section
	for i, line := range lines {
		if strings.HasPrefix(line, "flutter:") {
			// Flutter section already exists, this shouldn't happen
			return fmt.Errorf("flutter section exists but wasn't detected properly")
		}

		// Insert before any of these sections if they exist
		if strings.HasPrefix(line, "dev_dependencies:") ||
			strings.HasPrefix(line, "executables:") ||
			strings.HasPrefix(line, "platforms:") {
			insertIndex = i
			break
		}
	}

	// Create the flutter section with assets
	flutterSection := []string{
		"",
		"flutter:",
		"  assets:",
		"    - " + assetPath,
	}

	// Insert the flutter section
	updatedLines := append(
		lines[:insertIndex],
		append(flutterSection, lines[insertIndex:]...)...,
	)

	// Write the updated content back to the file
	updatedContent := strings.Join(updatedLines, "\n")
	return os.WriteFile(pubspecPath, []byte(updatedContent), 0644)
}

// updatePubspecWithNewAssetsSection adds a new assets section to the flutter section in pubspec.yaml
func updatePubspecWithNewAssetsSection(pubspecPath, pubspecContent, assetPath string) error {
	lines := strings.Split(pubspecContent, "\n")
	flutterLineIndex := -1
	flutterIndentation := ""

	// Find the flutter section
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "flutter:" {
			flutterLineIndex = i
			// Calculate the indentation of the flutter section
			flutterIndentation = strings.Repeat(" ", len(line)-len(trimmedLine))
			break
		}
	}

	if flutterLineIndex == -1 {
		// Flutter section not found, this shouldn't happen
		return fmt.Errorf("flutter section not found")
	}

	// Find where to insert the assets section (after the flutter line)
	insertIndex := flutterLineIndex + 1
	childIndentation := flutterIndentation + "  "

	// Create the assets section
	assetsSection := []string{
		childIndentation + "assets:",
		childIndentation + "  - " + assetPath,
	}

	// Insert the assets section
	updatedLines := append(
		lines[:insertIndex],
		append(assetsSection, lines[insertIndex:]...)...,
	)

	// Write the updated content back to the file
	updatedContent := strings.Join(updatedLines, "\n")
	return os.WriteFile(pubspecPath, []byte(updatedContent), 0644)
}

// updatePubspecWithNewAssetPath adds a new asset path to the existing assets section in pubspec.yaml
func updatePubspecWithNewAssetPath(pubspecPath, pubspecContent, assetPath string) error {
	lines := strings.Split(pubspecContent, "\n")
	assetsLineIndex := -1
	assetsIndentation := ""

	// Find the assets section
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "assets:" {
			assetsLineIndex = i
			// Calculate the indentation of the assets section
			assetsIndentation = strings.Repeat(" ", len(line)-len(trimmedLine))
			break
		}
	}

	if assetsLineIndex == -1 {
		// Assets section not found, this shouldn't happen
		return fmt.Errorf("assets section not found")
	}

	// Find the last asset entry to determine where to insert the new asset
	lastAssetIndex := assetsLineIndex
	assetEntryIndentation := assetsIndentation + "  "

	for i := assetsLineIndex + 1; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 {
			continue
		}

		// If this line has less indentation than expected for an asset entry,
		// we've moved out of the assets section
		if !strings.HasPrefix(line, assetEntryIndentation) {
			break
		}

		// If this line starts with "- ", it's an asset entry
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "- ") {
			lastAssetIndex = i
		}
	}

	// Insert the new asset entry after the last one
	insertIndex := lastAssetIndex + 1
	newAssetLine := assetEntryIndentation + "- " + assetPath

	// Insert the new asset line
	updatedLines := append(
		lines[:insertIndex],
		append([]string{newAssetLine}, lines[insertIndex:]...)...,
	)

	// Write the updated content back to the file
	updatedContent := strings.Join(updatedLines, "\n")
	return os.WriteFile(pubspecPath, []byte(updatedContent), 0644)
}

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

  /// Fonts assets
  static final Fonts fonts = Fonts._();

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

  /// Translation assets
  static final Translations translations = Translations._();
}

`)

	// Add asset classes
	addAssetClass(&content, "Images", ImageAsset, assets[ImageAsset])
	addAssetClass(&content, "Fonts", FontAsset, assets[FontAsset])
	addAssetClass(&content, "Animations", AnimationAsset, assets[AnimationAsset])
	addAssetClass(&content, "Audio", AudioAsset, assets[AudioAsset])
	addAssetClass(&content, "Videos", VideoAsset, assets[VideoAsset])
	addAssetClass(&content, "Json", JSONAsset, assets[JSONAsset])
	addAssetClass(&content, "Svgs", SVGAsset, assets[SVGAsset])
	addAssetClass(&content, "Translations", TranslationAsset, assets[TranslationAsset])

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
