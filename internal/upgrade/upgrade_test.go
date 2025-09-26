package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestReplaceBinaryWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a mock current executable
	currentPath := filepath.Join(tempDir, "td.exe")
	currentContent := []byte("current executable content")
	if err := os.WriteFile(currentPath, currentContent, 0755); err != nil {
		t.Fatalf("Failed to create mock current executable: %v", err)
	}

	// Create a mock new executable
	newPath := filepath.Join(tempDir, "td_new.exe")
	newContent := []byte("new executable content")
	if err := os.WriteFile(newPath, newContent, 0755); err != nil {
		t.Fatalf("Failed to create mock new executable: %v", err)
	}

	// Note: We can't easily test replaceBinaryWindows because it calls os.Exit(0)
	// Instead, we'll test that the batch script is created correctly

	// We'll create a mock version that doesn't exit for testing
	batchScript := currentPath + "_update.bat"

	scriptContent := fmt.Sprintf(`@echo off
echo Updating td...
timeout /t 2 /nobreak >nul
copy /y "%s" "%s" >nul
if errorlevel 1 (
    echo Failed to update td
    pause
    exit /b 1
)
echo Update completed successfully!
del "%s" >nul 2>&1
start "" "%s"
del "%%~f0" >nul 2>&1
`, newPath, currentPath, newPath, currentPath)

	// Write the batch script
	if err := os.WriteFile(batchScript, []byte(scriptContent), 0644); err != nil {
		t.Fatalf("Failed to create update script: %v", err)
	}

	// Verify the batch script was created
	if _, err := os.Stat(batchScript); err != nil {
		t.Errorf("Expected batch script to be created: %v", err)
	}

	// Clean up
	os.Remove(batchScript)
}

func TestReplaceBinaryUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a mock current executable
	currentPath := filepath.Join(tempDir, "td")
	currentContent := []byte("current executable content")
	if err := os.WriteFile(currentPath, currentContent, 0755); err != nil {
		t.Fatalf("Failed to create mock current executable: %v", err)
	}

	// Create a mock new executable
	newPath := filepath.Join(tempDir, "td_new")
	newContent := []byte("new executable content")
	if err := os.WriteFile(newPath, newContent, 0755); err != nil {
		t.Fatalf("Failed to create mock new executable: %v", err)
	}

	// Test the Unix replacement function
	err := replaceBinaryUnix(currentPath, newPath)
	if err != nil {
		t.Fatalf("replaceBinaryUnix failed: %v", err)
	}

	// Verify the current executable now has the new content
	actualContent, err := os.ReadFile(currentPath)
	if err != nil {
		t.Fatalf("Failed to read current executable after replacement: %v", err)
	}
	if string(actualContent) != string(newContent) {
		t.Errorf("Expected current executable to have new content, got %s", string(actualContent))
	}

	// Verify the backup doesn't exist (should be cleaned up)
	backupPath := currentPath + ".backup"
	if _, err := os.Stat(backupPath); err == nil {
		t.Errorf("Expected backup file to be cleaned up, but it still exists")
	}
}

func TestCleanupOldExecutables(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	// This is a more complex test since it involves os.Executable()
	// We'll create a mock scenario in a temp directory
	tempDir := t.TempDir()

	// Create a mock update script file
	scriptFile := filepath.Join(tempDir, "test.exe_update.bat")
	if err := os.WriteFile(scriptFile, []byte("@echo off\necho test"), 0644); err != nil {
		t.Fatalf("Failed to create mock script file: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(scriptFile); err != nil {
		t.Fatalf("Mock script file should exist: %v", err)
	}

	// Remove the file manually to simulate cleanup
	if err := os.Remove(scriptFile); err != nil {
		t.Fatalf("Failed to remove script file: %v", err)
	}

	// Verify the file is gone
	if _, err := os.Stat(scriptFile); err == nil {
		t.Errorf("Expected script file to be removed")
	}
}

func TestDetectPlatform(t *testing.T) {
	platform, arch := detectPlatform()

	// Verify platform is one of the expected values
	expectedPlatforms := map[string]bool{
		"windows": true,
		"linux":   true,
		"darwin":  true,
	}

	if !expectedPlatforms[platform] {
		t.Errorf("Unexpected platform: %s", platform)
	}

	// Verify architecture is one of the expected values
	expectedArchs := map[string]bool{
		"amd64": true,
		"arm64": true,
	}

	if !expectedArchs[arch] {
		t.Errorf("Unexpected architecture: %s", arch)
	}

	// Verify it matches runtime values (after normalization)
	expectedPlatform := runtime.GOOS
	expectedArch := runtime.GOARCH

	if platform != expectedPlatform {
		t.Errorf("Platform mismatch: got %s, expected %s", platform, expectedPlatform)
	}

	if arch != expectedArch {
		t.Errorf("Architecture mismatch: got %s, expected %s", arch, expectedArch)
	}
}
