package upgrade

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/voioo/td/internal/logger"
)

const (
	githubRepo = "voioo/td"
	baseURL    = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// Upgrade performs the self-upgrade process
func Upgrade(currentVersion string) error {
	logger.Info("Starting upgrade process", logger.F("current_version", currentVersion))

	// Get platform and architecture info
	platform, arch := detectPlatform()
	logger.Info("Detected platform",
		logger.F("platform", platform),
		logger.F("arch", arch))

	// Get latest release info
	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)
	}

	logger.Info("Latest release found",
		logger.F("version", release.TagName),
		logger.F("current_version", currentVersion))

	// Check if upgrade is needed
	if release.TagName == currentVersion {
		logger.Info("Already at latest version")
		fmt.Println("Already at latest version:", currentVersion)
		return nil
	}

	// Find appropriate asset
	asset, checksumAsset, err := findAsset(release, platform, arch)
	if err != nil {
		return fmt.Errorf("failed to find suitable asset: %w", err)
	}

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Check if installed via package manager
	if isPackageManaged() {
		return fmt.Errorf("td appears to be installed via a package manager. Please use your package manager to upgrade")
	}

	// Download and verify the binary
	tempDir, err := os.MkdirTemp("", "td-upgrade-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	logger.Info("Downloading new version",
		logger.F("asset", asset.Name),
		logger.F("temp_dir", tempDir))

	newBinaryPath, err := downloadAndExtract(asset, checksumAsset, tempDir)
	if err != nil {
		return fmt.Errorf("failed to download and extract: %w", err)
	}

	// Replace the current binary
	if err := replaceBinary(exePath, newBinaryPath); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	logger.Info("Upgrade completed successfully",
		logger.F("from", currentVersion),
		logger.F("to", release.TagName))

	return nil
}

// detectPlatform returns the platform and architecture strings
func detectPlatform() (platform, arch string) {
	platform = runtime.GOOS
	arch = runtime.GOARCH

	// Normalize platform names to match release naming
	switch platform {
	case "darwin":
		platform = "darwin"
	case "linux":
		platform = "linux"
	case "windows":
		platform = "windows"
	}

	// Normalize architecture names
	switch arch {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	}

	return platform, arch
}

// getLatestRelease fetches the latest release information from GitHub
func getLatestRelease() (*Release, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	// Set User-Agent to avoid rate limiting
	req.Header.Set("User-Agent", "td-upgrade/"+runtime.Version())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// findAsset finds the appropriate asset for the current platform/architecture
func findAsset(release *Release, platform, arch string) (binaryAsset, checksumAsset Asset, err error) {
	binaryName := fmt.Sprintf("td_%s_%s", platform, arch)
	checksumName := binaryName + ".sha256"

	var archiveExt string
	if platform == "windows" {
		archiveExt = ".zip"
	} else {
		archiveExt = ".tar.gz"
	}

	archiveName := binaryName + archiveExt
	checksumArchiveName := checksumName

	for _, asset := range release.Assets {
		if asset.Name == archiveName {
			binaryAsset = asset
		}
		if asset.Name == checksumArchiveName {
			checksumAsset = asset
		}
	}

	if binaryAsset.Name == "" {
		return Asset{}, Asset{}, fmt.Errorf("no suitable binary found for %s/%s", platform, arch)
	}

	if checksumAsset.Name == "" {
		return Asset{}, Asset{}, fmt.Errorf("no checksum file found for %s/%s", platform, arch)
	}

	return binaryAsset, checksumAsset, nil
}

// downloadAndExtract downloads the asset, verifies checksum, and extracts the binary
func downloadAndExtract(asset, checksumAsset Asset, tempDir string) (string, error) {
	// Download checksum
	checksumPath := filepath.Join(tempDir, checksumAsset.Name)
	if err := downloadFile(checksumAsset.BrowserDownloadURL, checksumPath); err != nil {
		return "", fmt.Errorf("failed to download checksum: %w", err)
	}

	// Read expected checksum
	checksumData, err := os.ReadFile(checksumPath)
	if err != nil {
		return "", fmt.Errorf("failed to read checksum file: %w", err)
	}

	expectedChecksum := strings.TrimSpace(strings.Split(string(checksumData), " ")[0])

	// Download archive
	archivePath := filepath.Join(tempDir, asset.Name)
	if err := downloadFile(asset.BrowserDownloadURL, archivePath); err != nil {
		return "", fmt.Errorf("failed to download archive: %w", err)
	}

	// Verify checksum
	if err := verifyChecksum(archivePath, expectedChecksum); err != nil {
		return "", fmt.Errorf("checksum verification failed: %w", err)
	}

	// Extract binary
	binaryPath, err := extractBinary(archivePath, tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to extract binary: %w", err)
	}

	return binaryPath, nil
}

// downloadFile downloads a file from the given URL to the specified path
func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// verifyChecksum verifies the SHA256 checksum of a file
func verifyChecksum(filePath, expectedChecksum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	actualChecksum := fmt.Sprintf("%x", hash.Sum(nil))
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// extractBinary extracts the td binary from the archive
func extractBinary(archivePath, destDir string) (string, error) {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	}
	return extractTarGz(archivePath, destDir)
}

// extractTarGz extracts a tar.gz archive
func extractTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Extract only the td binary
		if header.Name == "td" {
			destPath := filepath.Join(destDir, "td")
			outFile, err := os.Create(destPath)
			if err != nil {
				return "", err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return "", err
			}

			// Make executable
			if err := os.Chmod(destPath, 0755); err != nil {
				return "", err
			}

			return destPath, nil
		}
	}

	return "", fmt.Errorf("td binary not found in archive")
}

// extractZip extracts a zip archive
func extractZip(archivePath, destDir string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		// Extract only the td.exe binary
		if f.Name == "td.exe" {
			destPath := filepath.Join(destDir, "td.exe")
			outFile, err := os.Create(destPath)
			if err != nil {
				return "", err
			}
			defer outFile.Close()

			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			if _, err := io.Copy(outFile, rc); err != nil {
				return "", err
			}

			return destPath, nil
		}
	}

	return "", fmt.Errorf("td.exe binary not found in archive")
}

// replaceBinary replaces the current executable with the new one
func replaceBinary(currentPath, newPath string) error {
	if runtime.GOOS == "windows" {
		return replaceBinaryWindows(currentPath, newPath)
	}
	return replaceBinaryUnix(currentPath, newPath)
}

// replaceBinaryUnix replaces the binary on Unix-like systems
func replaceBinaryUnix(currentPath, newPath string) error {
	// Create backup of current binary
	backupPath := currentPath + ".backup"
	if err := copyFile(currentPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Replace the binary
	if err := os.Rename(newPath, currentPath); err != nil {
		// Try to restore backup on failure
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Clean up backup
	os.Remove(backupPath)

	return nil
}

// replaceBinaryWindows replaces the binary on Windows using a separate updater process
func replaceBinaryWindows(currentPath, newPath string) error {
	// On Windows, we cannot replace a running executable directly because the OS locks it.
	// We need to use a batch script that runs after this process exits.

	// Create a batch script that will handle the replacement
	batchScript := currentPath + "_update.bat"

	// Create the batch script content
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
		return fmt.Errorf("failed to create update script: %w", err)
	}

	logger.Info("Created update script", logger.F("script_path", batchScript))

	// Start the batch script
	cmd := exec.Command("cmd", "/c", batchScript)
	// Hide the command window on Windows
	setWindowsHidden(cmd)

	if err := cmd.Start(); err != nil {
		os.Remove(batchScript) // Clean up on failure
		return fmt.Errorf("failed to start update script: %w", err)
	}

	fmt.Println("Update process started. td will restart automatically...")
	logger.Info("Update script started, exiting current process")

	// Exit the current process to allow the update to proceed
	os.Exit(0)

	return nil // This line will never be reached, but needed for compilation
}

// copyFile copies a file from src to dst
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
	if err != nil {
		return err
	}

	return dstFile.Close()
}

// isPackageManaged checks if td is likely installed via a package manager
func isPackageManaged() bool {
	// Check for Homebrew on macOS
	if runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("brew"); err == nil {
			// Check if td is in Homebrew's Cellar
			if _, err := os.Stat("/opt/homebrew/Cellar/td-tui"); err == nil {
				return true
			}
			if _, err := os.Stat("/usr/local/Cellar/td-tui"); err == nil {
				return true
			}
		}
	}

	// Check for AUR on Linux
	if runtime.GOOS == "linux" {
		// Check if installed via pacman (Arch Linux)
		if _, err := exec.LookPath("pacman"); err == nil {
			cmd := exec.Command("pacman", "-Q", "td-tui")
			if err := cmd.Run(); err == nil {
				return true
			}
		}
	}

	return false
}

// CleanupOldExecutables removes old executable files left from previous upgrades
func CleanupOldExecutables() {
	if runtime.GOOS != "windows" {
		return // Only needed on Windows
	}

	exePath, err := os.Executable()
	if err != nil {
		return // Can't determine executable path
	}

	// Clean up any leftover update script
	batchScript := exePath + "_update.bat"
	if _, err := os.Stat(batchScript); err == nil {
		if err := os.Remove(batchScript); err != nil {
			logger.Info("Failed to cleanup update script",
				logger.F("script_path", batchScript),
				logger.F("error", err))
		} else {
			logger.Info("Cleaned up update script",
				logger.F("script_path", batchScript))
		}
	}
}
