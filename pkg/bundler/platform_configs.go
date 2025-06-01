package bundler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Android platform handlers

// detectAndroidGradleFormat detects whether the project uses .gradle or .gradle.kts
func detectAndroidGradleFormat(projectPath string) (string, error) {
	ktsPath := filepath.Join(projectPath, "android", "app", "build.gradle.kts")
	groovyPath := filepath.Join(projectPath, "android", "app", "build.gradle")

	if _, err := os.Stat(ktsPath); err == nil {
		return "kts", nil
	}
	if _, err := os.Stat(groovyPath); err == nil {
		return "groovy", nil
	}
	return "", fmt.Errorf("no Android build file found")
}

func getAndroidBundleID(projectPath string) BundleIDInfo {
	info := BundleIDInfo{Platform: PlatformAndroid, Available: true}

	// Detect format
	format, err := detectAndroidGradleFormat(projectPath)
	if err != nil {
		info.Error = err.Error()
		return info
	}

	// Read appropriate file
	var buildFilePath string
	if format == "kts" {
		buildFilePath = filepath.Join(projectPath, "android", "app", "build.gradle.kts")
	} else {
		buildFilePath = filepath.Join(projectPath, "android", "app", "build.gradle")
	}

	content, err := os.ReadFile(buildFilePath)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to read build file: %v", err)
		return info
	}

	// Parse based on format
	var applicationID, namespace string
	if format == "kts" {
		applicationID, namespace = parseKotlinDSL(string(content))
	} else {
		applicationID, namespace = parseGroovyDSL(string(content))
	}

	info.BundleID = applicationID
	info.Namespace = namespace
	return info
}

func setAndroidBundleID(projectPath, bundleID string) error {
	// Detect format
	format, err := detectAndroidGradleFormat(projectPath)
	if err != nil {
		return err
	}

	// Read appropriate file
	var buildFilePath string
	if format == "kts" {
		buildFilePath = filepath.Join(projectPath, "android", "app", "build.gradle.kts")
	} else {
		buildFilePath = filepath.Join(projectPath, "android", "app", "build.gradle")
	}

	content, err := os.ReadFile(buildFilePath)
	if err != nil {
		return fmt.Errorf("failed to read build file: %v", err)
	}

	// Update based on format
	var newContent string
	if format == "kts" {
		newContent = updateKotlinDSL(string(content), bundleID)
	} else {
		newContent = updateGroovyDSL(string(content), bundleID)
	}

	// Write back to file
	return os.WriteFile(buildFilePath, []byte(newContent), 0644)
}

// parseKotlinDSL parses Kotlin DSL (.gradle.kts) format
func parseKotlinDSL(content string) (applicationID, namespace string) {
	// applicationId = "com.example.app"
	appIDRegex := regexp.MustCompile(`applicationId\s*=\s*"([^"]+)"`)
	if matches := appIDRegex.FindStringSubmatch(content); len(matches) > 1 {
		applicationID = matches[1]
	}

	// namespace = "com.example.app"
	namespaceRegex := regexp.MustCompile(`namespace\s*=\s*"([^"]+)"`)
	if matches := namespaceRegex.FindStringSubmatch(content); len(matches) > 1 {
		namespace = matches[1]
	}

	return applicationID, namespace
}

// parseGroovyDSL parses Groovy DSL (.gradle) format
func parseGroovyDSL(content string) (applicationID, namespace string) {
	// applicationId "com.example.app" or applicationId 'com.example.app'
	appIDRegex := regexp.MustCompile(`applicationId\s+["']([^"']+)["']`)
	if matches := appIDRegex.FindStringSubmatch(content); len(matches) > 1 {
		applicationID = matches[1]
	}

	// namespace "com.example.app" or namespace 'com.example.app'
	namespaceRegex := regexp.MustCompile(`namespace\s+["']([^"']+)["']`)
	if matches := namespaceRegex.FindStringSubmatch(content); len(matches) > 1 {
		namespace = matches[1]
	}

	return applicationID, namespace
}

// updateKotlinDSL updates Kotlin DSL (.gradle.kts) format
func updateKotlinDSL(content, newBundleID string) string {
	// Update applicationId = "old" to applicationId = "new"
	appIDRegex := regexp.MustCompile(`(applicationId\s*=\s*)"[^"]+"`)
	content = appIDRegex.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newBundleID))

	// Update namespace = "old" to namespace = "new"
	namespaceRegex := regexp.MustCompile(`(namespace\s*=\s*)"[^"]+"`)
	content = namespaceRegex.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newBundleID))

	return content
}

// updateGroovyDSL updates Groovy DSL (.gradle) format
func updateGroovyDSL(content, newBundleID string) string {
	// Update applicationId "old" or applicationId 'old' to applicationId "new"
	appIDRegex := regexp.MustCompile(`(applicationId\s+)["'][^"']+["']`)
	content = appIDRegex.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newBundleID))

	// Update namespace "old" or namespace 'old' to namespace "new"
	namespaceRegex := regexp.MustCompile(`(namespace\s+)["'][^"']+["']`)
	content = namespaceRegex.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newBundleID))

	return content
}

// iOS platform handlers

func getIOSBundleID(projectPath string) BundleIDInfo {
	info := BundleIDInfo{Platform: PlatformIOS, Available: true}

	projectFile := filepath.Join(projectPath, "ios", "Runner.xcodeproj", "project.pbxproj")

	bundleID, err := parseIOSProjectFile(projectFile)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse iOS project file: %v", err)
		return info
	}

	info.BundleID = bundleID
	return info
}

func setIOSBundleID(projectPath, bundleID string) error {
	projectFile := filepath.Join(projectPath, "ios", "Runner.xcodeproj", "project.pbxproj")
	return updateIOSProjectFile(projectFile, bundleID)
}

func parseIOSProjectFile(projectFile string) (string, error) {
	content, err := os.ReadFile(projectFile)
	if err != nil {
		return "", err
	}

	// Look for PRODUCT_BUNDLE_IDENTIFIER = com.example.app;
	bundleIDRegex := regexp.MustCompile(`PRODUCT_BUNDLE_IDENTIFIER\s*=\s*([^;]+);`)
	matches := bundleIDRegex.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1]), nil
	}

	return "", fmt.Errorf("PRODUCT_BUNDLE_IDENTIFIER not found")
}

func updateIOSProjectFile(projectFile, bundleID string) error {
	content, err := os.ReadFile(projectFile)
	if err != nil {
		return err
	}

	// Update PRODUCT_BUNDLE_IDENTIFIER = old; to PRODUCT_BUNDLE_IDENTIFIER = new;
	bundleIDRegex := regexp.MustCompile(`(PRODUCT_BUNDLE_IDENTIFIER\s*=\s*)[^;]+(;)`)
	newContent := bundleIDRegex.ReplaceAllString(string(content), fmt.Sprintf("${1}%s${2}", bundleID))

	return os.WriteFile(projectFile, []byte(newContent), 0644)
}

// macOS platform handlers

func getMacOSBundleID(projectPath string) BundleIDInfo {
	info := BundleIDInfo{Platform: PlatformMacOS, Available: true}

	configPath := filepath.Join(projectPath, "macos", "Runner", "Configs", "AppInfo.xcconfig")

	bundleID, err := parseXCConfig(configPath, "PRODUCT_BUNDLE_IDENTIFIER")
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse AppInfo.xcconfig: %v", err)
		return info
	}

	info.BundleID = bundleID
	return info
}

func setMacOSBundleID(projectPath, bundleID string) error {
	configPath := filepath.Join(projectPath, "macos", "Runner", "Configs", "AppInfo.xcconfig")
	return updateXCConfig(configPath, "PRODUCT_BUNDLE_IDENTIFIER", bundleID)
}

func parseXCConfig(configPath, key string) (string, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, key+" = ") {
			return strings.TrimSpace(strings.TrimPrefix(line, key+" = ")), nil
		}
	}

	return "", fmt.Errorf("key %s not found", key)
}

func updateXCConfig(configPath, key, value string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), key+" = ") {
			// Replace the line
			lines = append(lines, fmt.Sprintf("%s = %s", key, value))
		} else {
			lines = append(lines, line)
		}
	}
	file.Close()

	// Write back to file
	content := strings.Join(lines, "\n")
	return os.WriteFile(configPath, []byte(content), 0644)
}

// Linux platform handlers

func getLinuxBundleID(projectPath string) BundleIDInfo {
	info := BundleIDInfo{Platform: PlatformLinux, Available: true}

	cmakePath := filepath.Join(projectPath, "linux", "CMakeLists.txt")

	binaryName, err := parseCMakeFile(cmakePath, "BINARY_NAME")
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse CMakeLists.txt: %v", err)
		return info
	}

	// For Linux, we use the binary name as the bundle ID
	info.BundleID = binaryName
	return info
}

func setLinuxBundleID(projectPath, bundleID string) error {
	cmakePath := filepath.Join(projectPath, "linux", "CMakeLists.txt")
	return updateCMakeFile(cmakePath, "BINARY_NAME", bundleID)
}

// Windows platform handlers

func getWindowsBundleID(projectPath string) BundleIDInfo {
	info := BundleIDInfo{Platform: PlatformWindows, Available: true}

	cmakePath := filepath.Join(projectPath, "windows", "CMakeLists.txt")

	binaryName, err := parseCMakeFile(cmakePath, "BINARY_NAME")
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse CMakeLists.txt: %v", err)
		return info
	}

	// For Windows, we use the binary name as the bundle ID
	info.BundleID = binaryName
	return info
}

func setWindowsBundleID(projectPath, bundleID string) error {
	cmakePath := filepath.Join(projectPath, "windows", "CMakeLists.txt")

	// Update both project name and BINARY_NAME
	if err := updateCMakeProjectName(cmakePath, bundleID); err != nil {
		return err
	}

	return updateCMakeFile(cmakePath, "BINARY_NAME", bundleID)
}

func parseCMakeFile(cmakePath, variable string) (string, error) {
	file, err := os.Open(cmakePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(fmt.Sprintf(`set\(%s\s+"([^"]+)"\)`, variable))

	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("variable %s not found", variable)
}

func updateCMakeFile(cmakePath, variable, value string) error {
	content, err := os.ReadFile(cmakePath)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(fmt.Sprintf(`(set\(%s\s+")[^"]+("\))`, variable))
	newContent := re.ReplaceAllString(string(content), fmt.Sprintf("${1}%s${2}", value))

	return os.WriteFile(cmakePath, []byte(newContent), 0644)
}

func updateCMakeProjectName(cmakePath, bundleID string) error {
	content, err := os.ReadFile(cmakePath)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(project\()[^)]+(\s+LANGUAGES\s+CXX\))`)
	newContent := re.ReplaceAllString(string(content), fmt.Sprintf("${1}%s${2}", bundleID))

	return os.WriteFile(cmakePath, []byte(newContent), 0644)
}

// Web platform handlers

type WebManifest struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	StartURL  string `json:"start_url,omitempty"`
	ID        string `json:"id,omitempty"` // This can serve as bundle ID for web
}

func getWebBundleID(projectPath string) BundleIDInfo {
	info := BundleIDInfo{Platform: PlatformWeb, Available: true}

	// Get bundle ID from manifest.json
	manifestPath := filepath.Join(projectPath, "web", "manifest.json")
	bundleID, err := parseWebManifest(manifestPath)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse manifest.json: %v", err)
		return info
	}

	info.BundleID = bundleID
	return info
}

func setWebBundleID(projectPath, bundleID string) error {
	// Update manifest.json
	manifestPath := filepath.Join(projectPath, "web", "manifest.json")
	return updateWebManifest(manifestPath, bundleID)
}

func parseWebManifest(manifestPath string) (string, error) {
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}

	var manifest WebManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return "", err
	}

	// Use ID field if available, otherwise use name
	if manifest.ID != "" {
		return manifest.ID, nil
	}
	return manifest.Name, nil
}

func updateWebManifest(manifestPath, bundleID string) error {
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	var manifest map[string]interface{}
	if err := json.Unmarshal(content, &manifest); err != nil {
		return err
	}

	// Set the ID field as the bundle ID
	manifest["id"] = bundleID

	newContent, err := json.MarshalIndent(manifest, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(manifestPath, newContent, 0644)
}

// Backup and restore functions

func createPlatformBackup(projectPath string, platform Platform, backupDir string) error {
	var filesToBackup []string

	switch platform {
	case PlatformAndroid:
		// Check which format exists and backup accordingly
		ktsPath := "android/app/build.gradle.kts"
		groovyPath := "android/app/build.gradle"
		if _, err := os.Stat(filepath.Join(projectPath, ktsPath)); err == nil {
			filesToBackup = []string{ktsPath}
		} else {
			filesToBackup = []string{groovyPath}
		}
	case PlatformIOS:
		filesToBackup = []string{"ios/Runner.xcodeproj/project.pbxproj"}
	case PlatformMacOS:
		filesToBackup = []string{"macos/Runner/Configs/AppInfo.xcconfig"}
	case PlatformLinux:
		filesToBackup = []string{"linux/CMakeLists.txt"}
	case PlatformWindows:
		filesToBackup = []string{"windows/CMakeLists.txt"}
	case PlatformWeb:
		filesToBackup = []string{"web/manifest.json"}
	}

	for _, file := range filesToBackup {
		srcPath := filepath.Join(projectPath, file)
		dstPath := filepath.Join(backupDir, file)

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

func restorePlatformBackup(projectPath string, platform Platform, backupDir string) {
	var filesToRestore []string

	switch platform {
	case PlatformAndroid:
		// Check which format exists and restore accordingly
		ktsPath := "android/app/build.gradle.kts"
		groovyPath := "android/app/build.gradle"
		if _, err := os.Stat(filepath.Join(projectPath, ktsPath)); err == nil {
			filesToRestore = []string{ktsPath}
		} else {
			filesToRestore = []string{groovyPath}
		}
	case PlatformIOS:
		filesToRestore = []string{"ios/Runner.xcodeproj/project.pbxproj"}
	case PlatformMacOS:
		filesToRestore = []string{"macos/Runner/Configs/AppInfo.xcconfig"}
	case PlatformLinux:
		filesToRestore = []string{"linux/CMakeLists.txt"}
	case PlatformWindows:
		filesToRestore = []string{"windows/CMakeLists.txt"}
	case PlatformWeb:
		filesToRestore = []string{"web/manifest.json"}
	}

	for _, file := range filesToRestore {
		srcPath := filepath.Join(backupDir, file)
		dstPath := filepath.Join(projectPath, file)
		copyFile(srcPath, dstPath) // Ignore errors during restore
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
