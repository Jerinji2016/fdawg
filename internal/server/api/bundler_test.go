package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jerinji2016/fdawg/pkg/flutter"
)

// helper to create a test Flutter Android project
func setupTestFlutterAndroidProject(t *testing.T, oldPackage string) string {
	tempDir := t.TempDir()

	// pubspec.yaml
	pubspec := "name: test_app\nversion: 1.0.0\nenvironment:\n  sdk: '>=3.0.0 <4.0.0'\n"
	if err := os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(pubspec), 0644); err != nil {
		t.Fatalf("failed to create pubspec.yaml: %v", err)
	}

	appDir := filepath.Join(tempDir, "android", "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create android/app: %v", err)
	}

	buildGroovy := `plugins {
    id 'com.android.application'
}
android {
    namespace '` + oldPackage + `'
    defaultConfig {
        applicationId '` + oldPackage + `'
    }
}
`
	if err := os.WriteFile(filepath.Join(appDir, "build.gradle"), []byte(buildGroovy), 0644); err != nil {
		t.Fatalf("failed to write build.gradle: %v", err)
	}

	pkgRelPath := filepath.Join(strings.Split(oldPackage, ".")...)
	srcDir := filepath.Join(appDir, "src", "main", "kotlin", pkgRelPath)
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("failed to create kotlin dir: %v", err)
	}

	mainActivityContent := `package ` + oldPackage + `

import io.flutter.embedding.android.FlutterActivity

class MainActivity: FlutterActivity() {
}
`
	if err := os.WriteFile(filepath.Join(srcDir, "MainActivity.kt"), []byte(mainActivityContent), 0644); err != nil {
		t.Fatalf("failed to write MainActivity.kt: %v", err)
	}

	return tempDir
}

func TestWebAPI_HandleSetBundleIDs_Android(t *testing.T) {
	oldPkg := "com.example.webold"
	newPkg := "com.example.webnew"
	projectDir := setupTestFlutterAndroidProject(t, oldPkg)

	validationResult := &flutter.ValidationResult{
		IsValid:     true,
		ProjectPath: projectDir,
	}

	api := NewBundlerAPI(validationResult)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	// Send POST /api/bundler/set with platform android
	reqBody := SetBundleIDsAPIRequest{
		Platforms: map[string]string{
			"android": newPkg,
		},
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/bundler/set", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify MainActivity.kt was moved
	newMainPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "webnew", "MainActivity.kt")
	content, err := os.ReadFile(newMainPath)
	if err != nil {
		t.Fatalf("MainActivity.kt not found at new path %s: %v", newMainPath, err)
	}
	if !strings.Contains(string(content), "package "+newPkg) {
		t.Errorf("MainActivity.kt does not have updated package:\n%s", string(content))
	}

	// Verify old path is gone
	oldMainDir := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "webold")
	if _, err := os.Stat(oldMainDir); !os.IsNotExist(err) {
		t.Errorf("old directory still exists at %s", oldMainDir)
	}
}
