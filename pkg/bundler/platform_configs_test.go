package bundler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper to create a minimal Flutter project with Android structure
func setupFlutterAndroidProject(t *testing.T, isKts bool, isJava bool, oldPackage string) string {
	tempDir := t.TempDir()

	// pubspec.yaml to qualify as Flutter project
	pubspec := "name: test_app\nversion: 1.0.0\nenvironment:\n  sdk: '>=3.0.0 <4.0.0'\n"
	if err := os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(pubspec), 0644); err != nil {
		t.Fatalf("failed to create pubspec.yaml: %v", err)
	}

	appDir := filepath.Join(tempDir, "android", "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create android/app: %v", err)
	}

	// build.gradle / build.gradle.kts
	if isKts {
		buildKts := `plugins {
    id("com.android.application")
    id("kotlin-android")
}
android {
    namespace = "` + oldPackage + `"
    defaultConfig {
        applicationId = "` + oldPackage + `"
        minSdk = 21
        targetSdk = 34
    }
}
`
		if err := os.WriteFile(filepath.Join(appDir, "build.gradle.kts"), []byte(buildKts), 0644); err != nil {
			t.Fatalf("failed to write build.gradle.kts: %v", err)
		}
	} else {
		buildGroovy := `plugins {
    id 'com.android.application'
}
android {
    namespace '` + oldPackage + `'
    defaultConfig {
        applicationId '` + oldPackage + `'
        minSdk 21
        targetSdk 34
    }
}
`
		if err := os.WriteFile(filepath.Join(appDir, "build.gradle"), []byte(buildGroovy), 0644); err != nil {
			t.Fatalf("failed to write build.gradle: %v", err)
		}
	}

	// AndroidManifest.xml
	mainDir := filepath.Join(appDir, "src", "main")
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		t.Fatalf("failed to create src/main: %v", err)
	}

	manifest := `<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="` + oldPackage + `">
    <application
        android:label="test_app">
        <activity
            android:name="` + oldPackage + `.MainActivity"
            android:exported="true">
        </activity>
    </application>
</manifest>`
	if err := os.WriteFile(filepath.Join(mainDir, "AndroidManifest.xml"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write AndroidManifest.xml: %v", err)
	}

	// Source directory
	pkgRelPath := filepath.Join(strings.Split(oldPackage, ".")...)
	var srcDir string
	if isJava {
		srcDir = filepath.Join(mainDir, "java", pkgRelPath)
		if err := os.MkdirAll(srcDir, 0755); err != nil {
			t.Fatalf("failed to create java dir: %v", err)
		}
		mainActivityContent := `package ` + oldPackage + `;

import io.flutter.embedding.android.FlutterActivity;

public class MainActivity extends FlutterActivity {
}
`
		if err := os.WriteFile(filepath.Join(srcDir, "MainActivity.java"), []byte(mainActivityContent), 0644); err != nil {
			t.Fatalf("failed to write MainActivity.java: %v", err)
		}
	} else {
		srcDir = filepath.Join(mainDir, "kotlin", pkgRelPath)
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
	}

	return tempDir
}

func TestSetAndroidBundleID_KotlinGroovy(t *testing.T) {
	oldPkg := "com.example.oldapp"
	newPkg := "com.example.newapp"
	projectDir := setupFlutterAndroidProject(t, false, false, oldPkg)

	// Add a subpackage file to verify deep files and imports
	subDir := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "oldapp", "services")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subDir: %v", err)
	}
	serviceContent := `package com.example.oldapp.services

import com.example.oldapp.MainActivity

class BackgroundService {
}
`
	if err := os.WriteFile(filepath.Join(subDir, "BackgroundService.kt"), []byte(serviceContent), 0644); err != nil {
		t.Fatalf("failed to create BackgroundService.kt: %v", err)
	}

	// Execute rename
	err := setAndroidBundleID(projectDir, newPkg)
	if err != nil {
		t.Fatalf("setAndroidBundleID failed: %v", err)
	}

	// 1. Verify build.gradle
	buildContent, err := os.ReadFile(filepath.Join(projectDir, "android", "app", "build.gradle"))
	if err != nil {
		t.Fatalf("failed to read build.gradle: %v", err)
	}
	if !strings.Contains(string(buildContent), `applicationId "`+newPkg+`"`) {
		t.Errorf("build.gradle missing updated applicationId, got:\n%s", string(buildContent))
	}
	if !strings.Contains(string(buildContent), `namespace "`+newPkg+`"`) {
		t.Errorf("build.gradle missing updated namespace, got:\n%s", string(buildContent))
	}

	// 2. Verify MainActivity.kt location and content
	newMainActivityPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "newapp", "MainActivity.kt")
	mainContent, err := os.ReadFile(newMainActivityPath)
	if err != nil {
		t.Fatalf("MainActivity.kt not found at %s: %v", newMainActivityPath, err)
	}
	if !strings.Contains(string(mainContent), "package "+newPkg) {
		t.Errorf("MainActivity.kt does not have updated package declaration:\n%s", string(mainContent))
	}

	// 3. Verify BackgroundService.kt location and content
	newServicePath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "newapp", "services", "BackgroundService.kt")
	serviceContentUpdated, err := os.ReadFile(newServicePath)
	if err != nil {
		t.Fatalf("BackgroundService.kt not found at %s: %v", newServicePath, err)
	}
	if !strings.Contains(string(serviceContentUpdated), "package "+newPkg+".services") {
		t.Errorf("BackgroundService.kt package not updated:\n%s", string(serviceContentUpdated))
	}
	if !strings.Contains(string(serviceContentUpdated), "import "+newPkg+".MainActivity") {
		t.Errorf("BackgroundService.kt import not updated:\n%s", string(serviceContentUpdated))
	}

	// 4. Verify old directories are gone
	oldDirPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "oldapp")
	if _, err := os.Stat(oldDirPath); !os.IsNotExist(err) {
		t.Errorf("old package directory still exists at %s", oldDirPath)
	}

	// 5. Verify AndroidManifest.xml
	manifestContent, err := os.ReadFile(filepath.Join(projectDir, "android", "app", "src", "main", "AndroidManifest.xml"))
	if err != nil {
		t.Fatalf("failed to read AndroidManifest.xml: %v", err)
	}
	if !strings.Contains(string(manifestContent), `package="`+newPkg+`"`) {
		t.Errorf("AndroidManifest.xml package attribute not updated:\n%s", string(manifestContent))
	}
	if !strings.Contains(string(manifestContent), `android:name="`+newPkg+`.MainActivity"`) {
		t.Errorf("AndroidManifest.xml activity android:name not updated:\n%s", string(manifestContent))
	}
}

func TestSetAndroidBundleID_KotlinKts(t *testing.T) {
	oldPkg := "com.example.oldapp"
	newPkg := "com.company.nested.app"
	projectDir := setupFlutterAndroidProject(t, true, false, oldPkg)

	err := setAndroidBundleID(projectDir, newPkg)
	if err != nil {
		t.Fatalf("setAndroidBundleID failed: %v", err)
	}

	buildContent, err := os.ReadFile(filepath.Join(projectDir, "android", "app", "build.gradle.kts"))
	if err != nil {
		t.Fatalf("failed to read build.gradle.kts: %v", err)
	}
	if !strings.Contains(string(buildContent), `applicationId = "`+newPkg+`"`) {
		t.Errorf("build.gradle.kts missing updated applicationId, got:\n%s", string(buildContent))
	}
	if !strings.Contains(string(buildContent), `namespace = "`+newPkg+`"`) {
		t.Errorf("build.gradle.kts missing updated namespace, got:\n%s", string(buildContent))
	}

	newMainPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "company", "nested", "app", "MainActivity.kt")
	mainContent, err := os.ReadFile(newMainPath)
	if err != nil {
		t.Fatalf("MainActivity.kt not found at %s: %v", newMainPath, err)
	}
	if !strings.Contains(string(mainContent), "package "+newPkg) {
		t.Errorf("MainActivity.kt package declaration incorrect:\n%s", string(mainContent))
	}
}

func TestSetAndroidBundleID_Java(t *testing.T) {
	oldPkg := "com.example.oldjava"
	newPkg := "com.example.newjava"
	projectDir := setupFlutterAndroidProject(t, false, true, oldPkg)

	err := setAndroidBundleID(projectDir, newPkg)
	if err != nil {
		t.Fatalf("setAndroidBundleID failed: %v", err)
	}

	newMainPath := filepath.Join(projectDir, "android", "app", "src", "main", "java", "com", "example", "newjava", "MainActivity.java")
	mainContent, err := os.ReadFile(newMainPath)
	if err != nil {
		t.Fatalf("MainActivity.java not found at %s: %v", newMainPath, err)
	}
	if !strings.Contains(string(mainContent), "package "+newPkg+";") {
		t.Errorf("MainActivity.java package declaration with semicolon incorrect:\n%s", string(mainContent))
	}

	oldMainDir := filepath.Join(projectDir, "android", "app", "src", "main", "java", "com", "example", "oldjava")
	if _, err := os.Stat(oldMainDir); !os.IsNotExist(err) {
		t.Errorf("old Java package directory still exists at %s", oldMainDir)
	}
}

func TestSetAndroidBundleID_DifferentSegmentCounts(t *testing.T) {
	// From 3 segments to 5 segments
	oldPkg := "com.example.three"
	newPkg := "org.mycorp.team.product.app"
	projectDir := setupFlutterAndroidProject(t, false, false, oldPkg)

	err := setAndroidBundleID(projectDir, newPkg)
	if err != nil {
		t.Fatalf("setAndroidBundleID failed (3->5 segments): %v", err)
	}

	newMainPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "org", "mycorp", "team", "product", "app", "MainActivity.kt")
	if _, err := os.Stat(newMainPath); err != nil {
		t.Fatalf("MainActivity.kt not found at 5-segment path: %v", err)
	}

	// Now from 5 segments down to 2 segments
	finalPkg := "com.two"
	err = setAndroidBundleID(projectDir, finalPkg)
	if err != nil {
		t.Fatalf("setAndroidBundleID failed (5->2 segments): %v", err)
	}

	twoMainPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "two", "MainActivity.kt")
	if _, err := os.Stat(twoMainPath); err != nil {
		t.Fatalf("MainActivity.kt not found at 2-segment path: %v", err)
	}

	// org/mycorp/team/product/app should be completely removed
	orgPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "org")
	if _, err := os.Stat(orgPath); !os.IsNotExist(err) {
		t.Errorf("org directory was not pruned after all files were moved: %s", orgPath)
	}
}

func TestSetAndroidBundleID_NestedPackageRenames(t *testing.T) {
	// Rename from com.example (parent) to com.example.nested (child)
	oldPkg := "com.example"
	newPkg := "com.example.nested"
	projectDir := setupFlutterAndroidProject(t, false, false, oldPkg)

	err := setAndroidBundleID(projectDir, newPkg)
	if err != nil {
		t.Fatalf("setAndroidBundleID failed (parent -> child): %v", err)
	}

	childMainPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "nested", "MainActivity.kt")
	if _, err := os.Stat(childMainPath); err != nil {
		t.Fatalf("child MainActivity.kt not found: %v", err)
	}

	// Rename back from com.example.nested (child) to com.example (parent)
	err = setAndroidBundleID(projectDir, oldPkg)
	if err != nil {
		t.Fatalf("setAndroidBundleID failed (child -> parent): %v", err)
	}

	parentMainPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "MainActivity.kt")
	if _, err := os.Stat(parentMainPath); err != nil {
		t.Fatalf("parent MainActivity.kt not found: %v", err)
	}
}

func TestSetAndroidBundleID_PreservesSiblingPackages(t *testing.T) {
	oldPkg := "com.example.oldapp"
	newPkg := "com.example.newapp"
	projectDir := setupFlutterAndroidProject(t, false, false, oldPkg)

	// Create sibling package com.example.other and unrelated com.unrelated
	siblingDir := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "other")
	if err := os.MkdirAll(siblingDir, 0755); err != nil {
		t.Fatalf("failed to create siblingDir: %v", err)
	}
	siblingFile := filepath.Join(siblingDir, "Sibling.kt")
	siblingContent := "package com.example.other\n\nclass Sibling {}\n"
	if err := os.WriteFile(siblingFile, []byte(siblingContent), 0644); err != nil {
		t.Fatalf("failed to write Sibling.kt: %v", err)
	}

	unrelatedDir := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "unrelated")
	if err := os.MkdirAll(unrelatedDir, 0755); err != nil {
		t.Fatalf("failed to create unrelatedDir: %v", err)
	}
	unrelatedFile := filepath.Join(unrelatedDir, "Unrelated.kt")
	unrelatedContent := "package com.unrelated\n\nclass Unrelated {}\n"
	if err := os.WriteFile(unrelatedFile, []byte(unrelatedContent), 0644); err != nil {
		t.Fatalf("failed to write Unrelated.kt: %v", err)
	}

	// Rename
	err := setAndroidBundleID(projectDir, newPkg)
	if err != nil {
		t.Fatalf("setAndroidBundleID failed: %v", err)
	}

	// Verify old package is deleted
	oldDir := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "oldapp")
	if _, err := os.Stat(oldDir); !os.IsNotExist(err) {
		t.Errorf("old package directory still exists: %s", oldDir)
	}

	// Verify sibling is untouched
	readSibling, err := os.ReadFile(siblingFile)
	if err != nil {
		t.Fatalf("sibling file missing: %v", err)
	}
	if string(readSibling) != siblingContent {
		t.Errorf("sibling file was modified:\n%s", string(readSibling))
	}

	// Verify unrelated is untouched
	readUnrelated, err := os.ReadFile(unrelatedFile)
	if err != nil {
		t.Fatalf("unrelated file missing: %v", err)
	}
	if string(readUnrelated) != unrelatedContent {
		t.Errorf("unrelated file was modified:\n%s", string(readUnrelated))
	}
}

func TestSetAndroidBundleID_FailureMissingMainActivity(t *testing.T) {
	oldPkg := "com.example.oldapp"
	newPkg := "com.example.newapp"
	projectDir := setupFlutterAndroidProject(t, false, false, oldPkg)

	// Remove MainActivity.kt
	mainActivityPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "oldapp", "MainActivity.kt")
	if err := os.Remove(mainActivityPath); err != nil {
		t.Fatalf("failed to remove MainActivity.kt: %v", err)
	}

	err := setAndroidBundleID(projectDir, newPkg)
	if err == nil {
		t.Fatalf("expected error when MainActivity is missing, got nil")
	}
	if !strings.Contains(err.Error(), "MainActivity not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestSetAndroidBundleID_FailureMissingSourceDir(t *testing.T) {
	oldPkg := "com.example.oldapp"
	newPkg := "com.example.newapp"
	projectDir := setupFlutterAndroidProject(t, false, false, oldPkg)

	// Remove entire android/app/src/main/kotlin directory
	kotlinDir := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin")
	if err := os.RemoveAll(kotlinDir); err != nil {
		t.Fatalf("failed to remove kotlin dir: %v", err)
	}

	err := setAndroidBundleID(projectDir, newPkg)
	if err == nil {
		t.Fatalf("expected error when source directory is missing, got nil")
	}
}

func TestSetBundleIDs_EndToEndWithBackupAndRestore(t *testing.T) {
	oldPkg := "com.example.oldapp"
	newPkg := "com.example.newapp"
	projectDir := setupFlutterAndroidProject(t, false, false, oldPkg)

	request := &BundleIDRequest{
		ProjectPath: projectDir,
		Platforms: map[Platform]string{
			PlatformAndroid: newPkg,
		},
	}

	// Successful run
	err := SetBundleIDs(request)
	if err != nil {
		t.Fatalf("SetBundleIDs failed: %v", err)
	}

	newMainPath := filepath.Join(projectDir, "android", "app", "src", "main", "kotlin", "com", "example", "newapp", "MainActivity.kt")
	if _, err := os.Stat(newMainPath); err != nil {
		t.Fatalf("MainActivity.kt not found after SetBundleIDs: %v", err)
	}
}
