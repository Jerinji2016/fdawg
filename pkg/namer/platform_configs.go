package namer

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Android platform handlers

type AndroidManifest struct {
	XMLName     xml.Name `xml:"manifest"`
	Application struct {
		Label string `xml:"label,attr"`
	} `xml:"application"`
}

func getAndroidAppName(projectPath string) AppNameInfo {
	info := AppNameInfo{Platform: PlatformAndroid, Available: true}
	
	manifestPath := filepath.Join(projectPath, "android", "app", "src", "main", "AndroidManifest.xml")
	
	file, err := os.Open(manifestPath)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to open AndroidManifest.xml: %v", err)
		return info
	}
	defer file.Close()

	var manifest AndroidManifest
	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&manifest); err != nil {
		info.Error = fmt.Sprintf("Failed to parse AndroidManifest.xml: %v", err)
		return info
	}

	info.DisplayName = manifest.Application.Label
	return info
}

func setAndroidAppName(projectPath, appName string) error {
	manifestPath := filepath.Join(projectPath, "android", "app", "src", "main", "AndroidManifest.xml")
	
	// Read the file content
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read AndroidManifest.xml: %v", err)
	}

	// Use regex to replace the android:label attribute
	re := regexp.MustCompile(`android:label="[^"]*"`)
	newContent := re.ReplaceAllString(string(content), fmt.Sprintf(`android:label="%s"`, appName))

	// Write back to file
	return os.WriteFile(manifestPath, []byte(newContent), 0644)
}

// iOS platform handlers

func getIOSAppName(projectPath string) AppNameInfo {
	info := AppNameInfo{Platform: PlatformIOS, Available: true}
	
	plistPath := filepath.Join(projectPath, "ios", "Runner", "Info.plist")
	
	// Read and parse the plist file
	displayName, bundleName, err := parseIOSPlist(plistPath)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse Info.plist: %v", err)
		return info
	}

	info.DisplayName = displayName
	info.InternalName = bundleName
	return info
}

func setIOSAppName(projectPath, appName string) error {
	plistPath := filepath.Join(projectPath, "ios", "Runner", "Info.plist")
	return updateIOSPlist(plistPath, appName)
}

func parseIOSPlist(plistPath string) (displayName, bundleName string, err error) {
	content, err := os.ReadFile(plistPath)
	if err != nil {
		return "", "", err
	}

	contentStr := string(content)
	
	// Extract CFBundleDisplayName
	displayNameRe := regexp.MustCompile(`<key>CFBundleDisplayName</key>\s*<string>([^<]*)</string>`)
	if matches := displayNameRe.FindStringSubmatch(contentStr); len(matches) > 1 {
		displayName = matches[1]
	}

	// Extract CFBundleName
	bundleNameRe := regexp.MustCompile(`<key>CFBundleName</key>\s*<string>([^<]*)</string>`)
	if matches := bundleNameRe.FindStringSubmatch(contentStr); len(matches) > 1 {
		bundleName = matches[1]
	}

	return displayName, bundleName, nil
}

func updateIOSPlist(plistPath, appName string) error {
	content, err := os.ReadFile(plistPath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Update CFBundleDisplayName
	displayNameRe := regexp.MustCompile(`(<key>CFBundleDisplayName</key>\s*<string>)[^<]*</string>`)
	contentStr = displayNameRe.ReplaceAllString(contentStr, fmt.Sprintf("${1}%s</string>", appName))

	// Update CFBundleName
	bundleNameRe := regexp.MustCompile(`(<key>CFBundleName</key>\s*<string>)[^<]*</string>`)
	contentStr = bundleNameRe.ReplaceAllString(contentStr, fmt.Sprintf("${1}%s</string>", appName))

	return os.WriteFile(plistPath, []byte(contentStr), 0644)
}

// macOS platform handlers

func getMacOSAppName(projectPath string) AppNameInfo {
	info := AppNameInfo{Platform: PlatformMacOS, Available: true}
	
	configPath := filepath.Join(projectPath, "macos", "Runner", "Configs", "AppInfo.xcconfig")
	
	productName, err := parseXCConfig(configPath, "PRODUCT_NAME")
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse AppInfo.xcconfig: %v", err)
		return info
	}

	info.DisplayName = productName
	return info
}

func setMacOSAppName(projectPath, appName string) error {
	configPath := filepath.Join(projectPath, "macos", "Runner", "Configs", "AppInfo.xcconfig")
	return updateXCConfig(configPath, "PRODUCT_NAME", appName)
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

func getLinuxAppName(projectPath string) AppNameInfo {
	info := AppNameInfo{Platform: PlatformLinux, Available: true}
	
	cmakePath := filepath.Join(projectPath, "linux", "CMakeLists.txt")
	
	binaryName, err := parseCMakeFile(cmakePath, "BINARY_NAME")
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse CMakeLists.txt: %v", err)
		return info
	}

	info.DisplayName = binaryName
	return info
}

func setLinuxAppName(projectPath, appName string) error {
	cmakePath := filepath.Join(projectPath, "linux", "CMakeLists.txt")
	return updateCMakeFile(cmakePath, "BINARY_NAME", appName)
}

// Windows platform handlers

func getWindowsAppName(projectPath string) AppNameInfo {
	info := AppNameInfo{Platform: PlatformWindows, Available: true}
	
	cmakePath := filepath.Join(projectPath, "windows", "CMakeLists.txt")
	
	binaryName, err := parseCMakeFile(cmakePath, "BINARY_NAME")
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse CMakeLists.txt: %v", err)
		return info
	}

	info.DisplayName = binaryName
	return info
}

func setWindowsAppName(projectPath, appName string) error {
	cmakePath := filepath.Join(projectPath, "windows", "CMakeLists.txt")
	
	// Update both project name and BINARY_NAME
	if err := updateCMakeProjectName(cmakePath, appName); err != nil {
		return err
	}
	
	return updateCMakeFile(cmakePath, "BINARY_NAME", appName)
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

func updateCMakeProjectName(cmakePath, appName string) error {
	content, err := os.ReadFile(cmakePath)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(project\()[^)]+(\s+LANGUAGES\s+CXX\))`)
	newContent := re.ReplaceAllString(string(content), fmt.Sprintf("${1}%s${2}", appName))

	return os.WriteFile(cmakePath, []byte(newContent), 0644)
}

// Web platform handlers

type WebManifest struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

func getWebAppName(projectPath string) AppNameInfo {
	info := AppNameInfo{Platform: PlatformWeb, Available: true}
	
	// Get name from manifest.json
	manifestPath := filepath.Join(projectPath, "web", "manifest.json")
	manifestName, err := parseWebManifest(manifestPath)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to parse manifest.json: %v", err)
		return info
	}

	info.DisplayName = manifestName
	return info
}

func setWebAppName(projectPath, appName string) error {
	// Update manifest.json
	manifestPath := filepath.Join(projectPath, "web", "manifest.json")
	if err := updateWebManifest(manifestPath, appName); err != nil {
		return err
	}

	// Update index.html
	indexPath := filepath.Join(projectPath, "web", "index.html")
	return updateWebIndex(indexPath, appName)
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

	return manifest.Name, nil
}

func updateWebManifest(manifestPath, appName string) error {
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	var manifest map[string]interface{}
	if err := json.Unmarshal(content, &manifest); err != nil {
		return err
	}

	manifest["name"] = appName
	manifest["short_name"] = appName

	newContent, err := json.MarshalIndent(manifest, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(manifestPath, newContent, 0644)
}

func updateWebIndex(indexPath, appName string) error {
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Update title tag
	titleRe := regexp.MustCompile(`<title>[^<]*</title>`)
	contentStr = titleRe.ReplaceAllString(contentStr, fmt.Sprintf("<title>%s</title>", appName))

	// Update apple-mobile-web-app-title meta tag
	metaRe := regexp.MustCompile(`(<meta name="apple-mobile-web-app-title" content=")[^"]*"`)
	contentStr = metaRe.ReplaceAllString(contentStr, fmt.Sprintf(`${1}%s"`, appName))

	return os.WriteFile(indexPath, []byte(contentStr), 0644)
}

// Backup and restore functions

func createPlatformBackup(projectPath string, platform Platform, backupDir string) error {
	var filesToBackup []string

	switch platform {
	case PlatformAndroid:
		filesToBackup = []string{"android/app/src/main/AndroidManifest.xml"}
	case PlatformIOS:
		filesToBackup = []string{"ios/Runner/Info.plist"}
	case PlatformMacOS:
		filesToBackup = []string{"macos/Runner/Configs/AppInfo.xcconfig"}
	case PlatformLinux:
		filesToBackup = []string{"linux/CMakeLists.txt"}
	case PlatformWindows:
		filesToBackup = []string{"windows/CMakeLists.txt"}
	case PlatformWeb:
		filesToBackup = []string{"web/manifest.json", "web/index.html"}
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
		filesToRestore = []string{"android/app/src/main/AndroidManifest.xml"}
	case PlatformIOS:
		filesToRestore = []string{"ios/Runner/Info.plist"}
	case PlatformMacOS:
		filesToRestore = []string{"macos/Runner/Configs/AppInfo.xcconfig"}
	case PlatformLinux:
		filesToRestore = []string{"linux/CMakeLists.txt"}
	case PlatformWindows:
		filesToRestore = []string{"windows/CMakeLists.txt"}
	case PlatformWeb:
		filesToRestore = []string{"web/manifest.json", "web/index.html"}
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
