package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// CurrentVersion is injected at build time via -ldflags.
// Default is "0.0.0" for local development.
var CurrentVersion = "0.0.0"

// GitHub repository info
const (
	GitHubOwner = "hoangtran1411"
	GitHubRepo  = "convert-vni-to-unicode"

	// HTTP client configuration
	httpTimeout     = 30 * time.Second
	downloadTimeout = 5 * time.Minute
	maxDownloadSize = 200 * 1024 * 1024 // 200MB max download size
)

// UpdateInfo holds information about available updates
type UpdateInfo struct {
	Available   bool   `json:"available"`
	CurrentVer  string `json:"currentVersion"`
	LatestVer   string `json:"latestVersion"`
	DownloadURL string `json:"downloadUrl"`
	ReleaseURL  string `json:"releaseUrl"`
}

// GitHubRelease represents a GitHub release API response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// GetCurrentVersion returns the current app version
func (a *App) GetCurrentVersion() string {
	return CurrentVersion
}

// CheckForUpdate checks GitHub for newer versions
func (a *App) CheckForUpdate() UpdateInfo {
	info := UpdateInfo{
		Available:  false,
		CurrentVer: CurrentVersion,
	}

	// Create HTTP client with timeout to prevent hanging
	client := &http.Client{Timeout: httpTimeout}

	// Call GitHub API
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwner, GitHubRepo)
	resp, err := client.Get(url)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to check update: %v", err)
		return info
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			runtime.LogErrorf(a.ctx, "Failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return info
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return info
	}

	info.LatestVer = release.TagName
	info.ReleaseURL = release.HTMLURL

	// Find Windows exe asset
	for _, asset := range release.Assets {
		if strings.HasSuffix(strings.ToLower(asset.Name), ".exe") {
			info.DownloadURL = asset.BrowserDownloadURL
			break
		}
	}

	// Compare versions
	if info.LatestVer != "" && CompareVersions(info.LatestVer, CurrentVersion) {
		info.Available = true
	}

	return info
}

// CompareVersions returns true if v1 is newer than v2
func CompareVersions(v1, v2 string) bool {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	for i := 0; i < 3; i++ {
		if parts1[i] > parts2[i] {
			return true
		}
		if parts1[i] < parts2[i] {
			return false
		}
	}
	return false
}

func parseVersion(v string) [3]int {
	var result [3]int
	parts := strings.Split(v, ".")
	for i := 0; i < len(parts) && i < 3; i++ {
		_, _ = fmt.Sscanf(parts[i], "%d", &result[i])
	}
	return result
}

// PerformUpdate downloads and installs the new version
func (a *App) PerformUpdate(downloadURL string) (bool, error) {
	if downloadURL == "" {
		return false, fmt.Errorf("no download URL provided")
	}

	exePath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve to absolute path safely
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return false, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Security: Validate the executable path doesn't contain shell-dangerous characters
	// This prevents command injection attacks
	if strings.ContainsAny(exePath, `"&|<>^`) {
		return false, fmt.Errorf("executable path contains unsafe characters")
	}

	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "vni_update.exe")

	runtime.EventsEmit(a.ctx, "updateProgress", "Downloading update...")

	// Create HTTP client with timeout for download
	client := &http.Client{Timeout: downloadTimeout}

	resp, err := client.Get(downloadURL)
	if err != nil {
		return false, fmt.Errorf("failed to download: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Error is non-critical during update
	}()

	// Validate response status
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Limit download size to prevent memory exhaustion attacks
	limitedReader := io.LimitReader(resp.Body, maxDownloadSize)

	out, err := os.Create(tempFile)
	if err != nil {
		return false, fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = io.Copy(out, limitedReader)
	if closeErr := out.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if err != nil {
		return false, fmt.Errorf("failed to save update: %w", err)
	}

	runtime.EventsEmit(a.ctx, "updateProgress", "Installing update...")

	// Create batch script to swap files and restart
	// Note: Paths are validated above to prevent command injection
	batchPath := filepath.Join(tempDir, "update_vni.bat")
	batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
del "%s"
move /y "%s" "%s"
start "" "%s"
del "%%~f0"
`, exePath, tempFile, exePath, exePath)

	// Use 0600 permission for security (owner read/write only)
	if err := os.WriteFile(batchPath, []byte(batchContent), 0600); err != nil {
		return false, fmt.Errorf("failed to create update script: %w", err)
	}

	cmd := exec.Command("cmd", "/c", "start", "/min", "", batchPath)
	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("failed to start update script: %w", err)
	}

	runtime.Quit(a.ctx)
	return true, nil
}
