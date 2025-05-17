package asset

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

	// Parse the file line by line to find the correct flutter section
	pubspecContent := string(pubspecData)
	lines := strings.Split(pubspecContent, "\n")

	// First, find the top-level flutter section (not under dependencies)
	flutterLineIndex := -1
	inDependenciesSection := false
	dependenciesIndentation := ""

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check if we're entering the dependencies section
		if trimmedLine == "dependencies:" {
			inDependenciesSection = true
			dependenciesIndentation = strings.Repeat(" ", len(line)-len(trimmedLine))
			continue
		}

		// Check if we're leaving the dependencies section
		if inDependenciesSection && len(line) > 0 && !strings.HasPrefix(line, dependenciesIndentation+"  ") && !strings.HasPrefix(line, " ") {
			inDependenciesSection = false
		}

		// Look for the flutter section, but only at the top level (not under dependencies)
		if !inDependenciesSection && trimmedLine == "flutter:" {
			flutterLineIndex = i
			break
		}
	}

	// If flutter section doesn't exist, add it
	if flutterLineIndex == -1 {
		return addNewFlutterSection(pubspecPath, pubspecContent, assetPath)
	}

	// Now check if the assets section exists under the flutter section
	assetsLineIndex := -1
	flutterIndentation := ""

	// Calculate the indentation of the flutter section
	for i, char := range lines[flutterLineIndex] {
		if char != ' ' {
			flutterIndentation = strings.Repeat(" ", i)
			break
		}
	}

	// Expected indentation for flutter children
	childIndentation := flutterIndentation + "  "

	// Look for the assets section within the flutter section
	for i := flutterLineIndex + 1; i < len(lines); i++ {
		line := lines[i]

		// If we've reached the end of the flutter section, break
		if len(line) > 0 && !strings.HasPrefix(line, " ") {
			break
		}

		// If we've reached another top-level section, break
		if len(line) > 0 && !strings.HasPrefix(line, childIndentation) && strings.TrimSpace(line) != "" {
			break
		}

		// Check if this is the assets section
		if strings.HasPrefix(line, childIndentation) && strings.TrimSpace(line) == "assets:" {
			assetsLineIndex = i
			break
		}
	}

	// If assets section doesn't exist, add it
	if assetsLineIndex == -1 {
		return addAssetsSection(pubspecPath, pubspecContent, flutterLineIndex, childIndentation, assetPath)
	}

	// Now check if the specific asset path already exists
	assetEntryIndentation := childIndentation + "  "
	assetExists := false

	for i := assetsLineIndex + 1; i < len(lines); i++ {
		line := lines[i]

		// If we've reached the end of the assets section, break
		if len(line) > 0 && !strings.HasPrefix(line, assetEntryIndentation) && strings.TrimSpace(line) != "" {
			break
		}

		// Check if this is the asset entry we're looking for
		if strings.HasPrefix(line, assetEntryIndentation) && strings.TrimSpace(line) == "- "+assetPath {
			assetExists = true
			break
		}
	}

	// If the asset doesn't exist, add it
	if !assetExists {
		return addAssetEntry(pubspecPath, pubspecContent, assetsLineIndex, assetEntryIndentation, assetPath)
	}

	// Asset already exists, no need to update
	return nil
}

// addNewFlutterSection adds a new flutter section with assets to the pubspec.yaml file
func addNewFlutterSection(pubspecPath, pubspecContent, assetPath string) error {
	lines := strings.Split(pubspecContent, "\n")
	insertIndex := len(lines)

	// Look for a good place to insert the flutter section
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Insert before dev_dependencies if it exists
		if trimmedLine == "dev_dependencies:" {
			insertIndex = i
			break
		}
	}

	// Create the flutter section with assets
	flutterSection := []string{
		"",
		"flutter:",
		"  uses-material-design: true",
		"",
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

// addAssetsSection adds an assets section to the flutter section
func addAssetsSection(pubspecPath, pubspecContent string, flutterLineIndex int, childIndentation, assetPath string) error {
	lines := strings.Split(pubspecContent, "\n")

	// Find where to insert the assets section
	insertIndex := flutterLineIndex + 1

	// Look for a good place to insert within the flutter section
	for i := flutterLineIndex + 1; i < len(lines); i++ {
		line := lines[i]

		// If we've reached the end of the flutter section, break
		if len(line) > 0 && !strings.HasPrefix(line, childIndentation) && strings.TrimSpace(line) != "" {
			insertIndex = i
			break
		}

		// Update the insert index as we go
		insertIndex = i + 1
	}

	// Create the assets section
	assetsSection := []string{
		"",
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

// addAssetEntry adds a new asset entry to the assets section
func addAssetEntry(pubspecPath, pubspecContent string, assetsLineIndex int, assetEntryIndentation, assetPath string) error {
	lines := strings.Split(pubspecContent, "\n")

	// Find where to insert the asset entry
	insertIndex := assetsLineIndex + 1

	// Look for the last asset entry
	for i := assetsLineIndex + 1; i < len(lines); i++ {
		line := lines[i]

		// If we've reached the end of the assets section, break
		if len(line) > 0 && !strings.HasPrefix(line, assetEntryIndentation) && strings.TrimSpace(line) != "" {
			insertIndex = i
			break
		}

		// If this is an asset entry, update the insert index
		if strings.HasPrefix(line, assetEntryIndentation) && strings.HasPrefix(strings.TrimSpace(line), "- ") {
			insertIndex = i + 1
		}
	}

	// Create the asset entry
	assetEntry := assetEntryIndentation + "- " + assetPath

	// Insert the asset entry
	updatedLines := append(
		lines[:insertIndex],
		append([]string{assetEntry}, lines[insertIndex:]...)...,
	)

	// Write the updated content back to the file
	updatedContent := strings.Join(updatedLines, "\n")
	return os.WriteFile(pubspecPath, []byte(updatedContent), 0644)
}
