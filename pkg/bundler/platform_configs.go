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
	if info.BundleID == "" && namespace != "" {
		info.BundleID = namespace
	}
	info.Namespace = namespace
	return info
}

type mainActivityInfo struct {
	path        string
	filename    string
	sourceRoot  string
	packageName string
}

func setAndroidBundleID(projectPath, bundleID string) error {
	// Detect format
	format, err := detectAndroidGradleFormat(projectPath)
	if err != nil {
		return err
	}

	// Read appropriate build file
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

	// Discover Android source roots in Flutter android/app/src
	sourceRoots, err := findAndroidSourceRoots(projectPath)
	if err != nil {
		return err
	}

	// Locate MainActivity (.kt or .java)
	mainActivity, err := findMainActivity(sourceRoots)
	if err != nil {
		return err
	}

	oldSourcePackage := mainActivity.packageName
	if oldSourcePackage == "" {
		return fmt.Errorf("no package declaration found in MainActivity (%s)", mainActivity.path)
	}

	// Verify the expected source directory for the old package exists
	oldRelPath := filepath.Join(strings.Split(oldSourcePackage, ".")...)
	expectedOldDir := filepath.Join(mainActivity.sourceRoot, oldRelPath)
	if _, err := os.Stat(expectedOldDir); os.IsNotExist(err) {
		return fmt.Errorf("expected Android source directory not found: %s", expectedOldDir)
	}

	// Parse build file values for additional package references
	var oldAppID, oldNamespace string
	if format == "kts" {
		oldAppID, oldNamespace = parseKotlinDSL(string(content))
	} else {
		oldAppID, oldNamespace = parseGroovyDSL(string(content))
	}

	// Update Gradle build file
	var newContent string
	if format == "kts" {
		newContent = updateKotlinDSL(string(content), bundleID)
	} else {
		newContent = updateGroovyDSL(string(content), bundleID)
	}

	if err := os.WriteFile(buildFilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write build file: %v", err)
	}

	// Update AndroidManifest.xml files with all known old package identifiers
	packagesToUpdate := []string{oldSourcePackage}
	if oldAppID != "" && oldAppID != oldSourcePackage {
		packagesToUpdate = append(packagesToUpdate, oldAppID)
	}
	if oldNamespace != "" && oldNamespace != oldSourcePackage && oldNamespace != oldAppID {
		packagesToUpdate = append(packagesToUpdate, oldNamespace)
	}
	for _, pkg := range packagesToUpdate {
		updateAndroidManifestFiles(projectPath, pkg, bundleID)
	}

	// Relocate source directories and update source file contents if package changed
	if oldSourcePackage != bundleID {
		for _, root := range sourceRoots {
			if err := relocatePackageDir(root, oldSourcePackage, bundleID); err != nil {
				return fmt.Errorf("failed to relocate package directory in %s: %v", root, err)
			}
		}

		// Update package declarations and imports across all source roots
		for _, root := range sourceRoots {
			if err := updateSourceFiles(root, oldSourcePackage, bundleID); err != nil {
				return fmt.Errorf("failed to update source files in %s: %v", root, err)
			}
		}
	}

	// Verify MainActivity exists at the new path
	newRelPath := filepath.Join(strings.Split(bundleID, ".")...)
	expectedNewMainActivity := filepath.Join(mainActivity.sourceRoot, newRelPath, mainActivity.filename)
	if _, err := os.Stat(expectedNewMainActivity); os.IsNotExist(err) {
		return fmt.Errorf("MainActivity verification failed: not found at %s", expectedNewMainActivity)
	}

	return nil
}

// findAndroidSourceRoots finds Kotlin and Java source roots in Flutter's android/app/src
func findAndroidSourceRoots(projectPath string) ([]string, error) {
	srcDir := filepath.Join(projectPath, "android", "app", "src")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("Android source directory not found: %s", srcDir)
	}

	var roots []string
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read Android src directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		for _, lang := range []string{"kotlin", "java"} {
			candidate := filepath.Join(srcDir, entry.Name(), lang)
			if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
				roots = append(roots, candidate)
			}
		}
	}

	// Also check direct android/app/src/kotlin and android/app/src/java if present
	for _, lang := range []string{"kotlin", "java"} {
		candidate := filepath.Join(srcDir, lang)
		if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
			roots = append(roots, candidate)
		}
	}

	if len(roots) == 0 {
		return nil, fmt.Errorf("no Kotlin or Java source directory found under %s", srcDir)
	}

	return roots, nil
}

// findMainActivity searches for MainActivity.kt or MainActivity.java in source roots
func findMainActivity(sourceRoots []string) (*mainActivityInfo, error) {
	for _, root := range sourceRoots {
		var found *mainActivityInfo
		err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !info.IsDir() && (info.Name() == "MainActivity.kt" || info.Name() == "MainActivity.java") {
				content, readErr := os.ReadFile(path)
				if readErr != nil {
					return readErr
				}
				pkg := extractPackageDeclaration(string(content))
				found = &mainActivityInfo{
					path:        path,
					filename:    info.Name(),
					sourceRoot:  root,
					packageName: pkg,
				}
				return filepath.SkipAll
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if found != nil {
			return found, nil
		}
	}
	return nil, fmt.Errorf("MainActivity not found (expected MainActivity.kt or MainActivity.java)")
}

// extractPackageDeclaration parses package declaration from Kotlin or Java source
func extractPackageDeclaration(content string) string {
	re := regexp.MustCompile(`(?m)^\s*package\s+([a-zA-Z0-9_.]+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// relocatePackageDir moves files from old package directory to new package directory
func relocatePackageDir(sourceRoot, oldPackage, newPackage string) error {
	oldRelPath := filepath.Join(strings.Split(oldPackage, ".")...)
	oldFullDir := filepath.Join(sourceRoot, oldRelPath)

	if _, err := os.Stat(oldFullDir); os.IsNotExist(err) {
		// Old package dir not present in this source root
		return nil
	}

	newRelPath := filepath.Join(strings.Split(newPackage, ".")...)
	newFullDir := filepath.Join(sourceRoot, newRelPath)

	// Use a temporary directory in sourceRoot to handle overlapping/nested paths safely
	tempDir, err := os.MkdirTemp(sourceRoot, ".fdawg-rename-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Move all entries from oldFullDir to tempDir
	entries, err := os.ReadDir(oldFullDir)
	if err != nil {
		return fmt.Errorf("failed to read old package directory: %v", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(oldFullDir, entry.Name())
		dstPath := filepath.Join(tempDir, entry.Name())
		if err := moveEntry(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to move %s to temporary directory: %v", entry.Name(), err)
		}
	}

	// Remove oldFullDir (now empty)
	if err := os.Remove(oldFullDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove old package directory: %v", err)
	}

	// Prune empty parent directories up to sourceRoot
	pruneEmptyParentDirs(oldFullDir, sourceRoot)

	// Create newFullDir
	if err := os.MkdirAll(newFullDir, 0755); err != nil {
		return fmt.Errorf("failed to create new package directory: %v", err)
	}

	// Move all entries from tempDir to newFullDir
	tempEntries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temporary directory: %v", err)
	}

	for _, entry := range tempEntries {
		srcPath := filepath.Join(tempDir, entry.Name())
		dstPath := filepath.Join(newFullDir, entry.Name())
		if err := moveEntry(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to move %s to new package directory: %v", entry.Name(), err)
		}
	}

	return nil
}

// moveEntry moves a file or directory, with fallback to copy and delete
func moveEntry(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		if err := copyDir(src, dst); err != nil {
			return err
		}
		return os.RemoveAll(src)
	}

	if err := copyFile(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

// pruneEmptyParentDirs removes empty parent directories from dir up to (exclusive) stopAt
func pruneEmptyParentDirs(dir, stopAt string) {
	cleanStopAt := filepath.Clean(stopAt)
	current := filepath.Dir(dir)
	for {
		cleanCurrent := filepath.Clean(current)
		if cleanCurrent == cleanStopAt || !strings.HasPrefix(cleanCurrent, cleanStopAt) {
			break
		}
		entries, err := os.ReadDir(cleanCurrent)
		if err != nil || len(entries) > 0 {
			break
		}
		if err := os.Remove(cleanCurrent); err != nil {
			break
		}
		current = filepath.Dir(cleanCurrent)
	}
}

// updateSourceFiles updates package declarations and imports in all Kotlin/Java files under root
func updateSourceFiles(sourceRoot, oldPackage, newPackage string) error {
	return filepath.Walk(sourceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".fdawg") {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".kt" || ext == ".java" {
			return updateSourceFileContent(path, oldPackage, newPackage)
		}
		return nil
	})
}

// updateSourceFileContent updates package declaration and imports for a single source file
func updateSourceFileContent(filePath, oldPackage, newPackage string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	s := string(content)
	orig := s

	// Update package declarations:
	// Matches: package <oldPackage>
	//          package <oldPackage>;
	//          package <oldPackage>.<subpackage>
	//          package <oldPackage>.<subpackage>;
	pkgRegex := regexp.MustCompile(`(?m)^(\s*package\s+)` + regexp.QuoteMeta(oldPackage) + `([;.\s].*|$)`)
	s = pkgRegex.ReplaceAllString(s, `${1}`+newPackage+`${2}`)

	// Update import statements:
	// Matches: import <oldPackage>.<something>
	importRegex := regexp.MustCompile(`(?m)^(\s*import\s+)` + regexp.QuoteMeta(oldPackage) + `([;.\s].*|$)`)
	s = importRegex.ReplaceAllString(s, `${1}`+newPackage+`${2}`)

	if s != orig {
		return os.WriteFile(filePath, []byte(s), 0644)
	}
	return nil
}

// updateAndroidManifestFiles updates package and activity references in all AndroidManifest.xml files
func updateAndroidManifestFiles(projectPath, oldPackage, newPackage string) {
	srcDir := filepath.Join(projectPath, "android", "app", "src")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return
	}

	_ = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == "AndroidManifest.xml" {
			_ = updateAndroidManifest(path, oldPackage, newPackage)
		}
		return nil
	})
}

// updateAndroidManifest updates package and activity references in a single AndroidManifest.xml
func updateAndroidManifest(manifestPath, oldPackage, newPackage string) error {
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	s := string(content)
	orig := s

	// Update package="..." or package='...'
	pkgRegex := regexp.MustCompile(`(package\s*=\s*["'])` + regexp.QuoteMeta(oldPackage) + `(["'])`)
	s = pkgRegex.ReplaceAllString(s, fmt.Sprintf(`${1}%s${2}`, newPackage))

	// Update android:name="oldPackage.MainActivity" or similar
	activityRegex := regexp.MustCompile(`(android:name\s*=\s*["'])` + regexp.QuoteMeta(oldPackage) + `(\.[^"']*["'])`)
	s = activityRegex.ReplaceAllString(s, fmt.Sprintf(`${1}%s${2}`, newPackage))

	if s != orig {
		return os.WriteFile(manifestPath, []byte(s), 0644)
	}
	return nil
}

// parseKotlinDSL parses Kotlin DSL (.gradle.kts) format
func parseKotlinDSL(content string) (applicationID, namespace string) {
	// applicationId = "com.example.app"
	appIDRegex := regexp.MustCompile(`applicationId\s*=\s*["']([^"']+)["']`)
	if matches := appIDRegex.FindStringSubmatch(content); len(matches) > 1 {
		applicationID = matches[1]
	}

	// namespace = "com.example.app"
	namespaceRegex := regexp.MustCompile(`namespace\s*=\s*["']([^"']+)["']`)
	if matches := namespaceRegex.FindStringSubmatch(content); len(matches) > 1 {
		namespace = matches[1]
	}

	return applicationID, namespace
}

// parseGroovyDSL parses Groovy DSL (.gradle) format
func parseGroovyDSL(content string) (applicationID, namespace string) {
	// applicationId "com.example.app" or applicationId = "com.example.app"
	appIDRegex := regexp.MustCompile(`applicationId(?:\s*=\s*|\s+)["']([^"']+)["']`)
	if matches := appIDRegex.FindStringSubmatch(content); len(matches) > 1 {
		applicationID = matches[1]
	}

	// namespace "com.example.app" or namespace = "com.example.app"
	namespaceRegex := regexp.MustCompile(`namespace(?:\s*=\s*|\s+)["']([^"']+)["']`)
	if matches := namespaceRegex.FindStringSubmatch(content); len(matches) > 1 {
		namespace = matches[1]
	}

	return applicationID, namespace
}

// updateKotlinDSL updates Kotlin DSL (.gradle.kts) format
func updateKotlinDSL(content, newBundleID string) string {
	// Update applicationId = "old" to applicationId = "new"
	appIDRegex := regexp.MustCompile(`(applicationId\s*=\s*)["'][^"']+["']`)
	content = appIDRegex.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newBundleID))

	// Update namespace = "old" to namespace = "new"
	namespaceRegex := regexp.MustCompile(`(namespace\s*=\s*)["'][^"']+["']`)
	content = namespaceRegex.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newBundleID))

	return content
}

// updateGroovyDSL updates Groovy DSL (.gradle) format
func updateGroovyDSL(content, newBundleID string) string {
	// Update applicationId "old" or applicationId = "old" to applicationId "new"
	appIDRegex := regexp.MustCompile(`(applicationId(?:\s*=\s*|\s+))["'][^"']+["']`)
	content = appIDRegex.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newBundleID))

	// Update namespace "old" or namespace = "old" to namespace "new"
	namespaceRegex := regexp.MustCompile(`(namespace(?:\s*=\s*|\s+))["'][^"']+["']`)
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

		// Also backup android/app/src directory if it exists
		srcDir := filepath.Join(projectPath, "android", "app", "src")
		if fi, err := os.Stat(srcDir); err == nil && fi.IsDir() {
			dstSrcDir := filepath.Join(backupDir, "android", "app", "src")
			_ = copyDir(srcDir, dstSrcDir)
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

		// Restore android/app/src directory if backed up
		backupSrcDir := filepath.Join(backupDir, "android", "app", "src")
		if fi, err := os.Stat(backupSrcDir); err == nil && fi.IsDir() {
			dstSrcDir := filepath.Join(projectPath, "android", "app", "src")
			_ = os.RemoveAll(dstSrcDir)
			_ = copyDir(backupSrcDir, dstSrcDir)
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

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
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
