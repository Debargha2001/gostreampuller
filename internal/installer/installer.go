package installer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ytdlpVersion  = "2024.11.18"
	ffmpegVersion = "7.1"
)

// GetBinariesDir returns the directory where binaries are stored
func GetBinariesDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	binDir := filepath.Join(homeDir, ".gostreampuller", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create bin directory: %w", err)
	}
	
	return binDir, nil
}

// GetYTDLPPath returns the path to locally installed yt-dlp
func GetYTDLPPath() (string, error) {
	binDir, err := GetBinariesDir()
	if err != nil {
		return "", err
	}
	
	executable := "yt-dlp"
	if runtime.GOOS == "windows" {
		executable = "yt-dlp.exe"
	}
	
	path := filepath.Join(binDir, executable)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	
	return "", fmt.Errorf("yt-dlp not found at %s", path)
}

// GetFFMPEGPath returns the path to locally installed ffmpeg
func GetFFMPEGPath() (string, error) {
	binDir, err := GetBinariesDir()
	if err != nil {
		return "", err
	}
	
	executable := "ffmpeg"
	if runtime.GOOS == "windows" {
		executable = "ffmpeg.exe"
	}
	
	path := filepath.Join(binDir, executable)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	
	return "", fmt.Errorf("ffmpeg not found at %s", path)
}

// InstallYTDLP downloads and installs yt-dlp
func InstallYTDLP(progressFn func(string)) error {
	binDir, err := GetBinariesDir()
	if err != nil {
		return err
	}
	
	var downloadURL string
	var executable string
	
	switch runtime.GOOS {
	case "linux":
		downloadURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
		executable = "yt-dlp"
	case "darwin":
		downloadURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"
		executable = "yt-dlp"
	case "windows":
		downloadURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
		executable = "yt-dlp.exe"
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	
	destPath := filepath.Join(binDir, executable)
	
	if progressFn != nil {
		progressFn(fmt.Sprintf("Downloading yt-dlp from %s...", downloadURL))
	}
	
	if err := downloadFile(downloadURL, destPath, progressFn); err != nil {
		return fmt.Errorf("failed to download yt-dlp: %w", err)
	}
	
	// Make executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(destPath, 0755); err != nil {
			return fmt.Errorf("failed to make yt-dlp executable: %w", err)
		}
	}
	
	if progressFn != nil {
		progressFn(fmt.Sprintf("✓ yt-dlp installed at: %s", destPath))
	}
	
	return nil
}

// InstallFFMPEG downloads and installs ffmpeg
func InstallFFMPEG(progressFn func(string)) error {
	binDir, err := GetBinariesDir()
	if err != nil {
		return err
	}
	
	var downloadURL string
	var needsExtraction bool
	var archiveType string // "zip" or "tar.gz"
	
	switch runtime.GOOS {
	case "linux":
		// Use static build from johnvansickle
		downloadURL = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
		needsExtraction = true
		archiveType = "tar.xz"
	case "darwin":
		// For macOS, recommend using brew but provide alternative
		if progressFn != nil {
			progressFn("For macOS, we recommend installing via Homebrew: brew install ffmpeg")
			progressFn("Attempting to download pre-built binary...")
		}
		downloadURL = "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip"
		needsExtraction = true
		archiveType = "zip"
	case "windows":
		downloadURL = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
		needsExtraction = true
		archiveType = "zip"
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	
	if progressFn != nil {
		progressFn(fmt.Sprintf("Downloading ffmpeg from %s...", downloadURL))
	}
	
	if needsExtraction {
		// Download to temp file
		tmpFile := filepath.Join(os.TempDir(), "ffmpeg-download")
		if err := downloadFile(downloadURL, tmpFile, progressFn); err != nil {
			return fmt.Errorf("failed to download ffmpeg: %w", err)
		}
		defer os.Remove(tmpFile)
		
		if progressFn != nil {
			progressFn("Extracting ffmpeg...")
		}
		
		// Extract based on archive type
		if archiveType == "zip" {
			if err := extractFFMPEGFromZip(tmpFile, binDir, progressFn); err != nil {
				return fmt.Errorf("failed to extract ffmpeg: %w", err)
			}
		} else if archiveType == "tar.xz" || archiveType == "tar.gz" {
			if err := extractFFMPEGFromTar(tmpFile, binDir, progressFn); err != nil {
				return fmt.Errorf("failed to extract ffmpeg: %w", err)
			}
		}
	} else {
		executable := "ffmpeg"
		if runtime.GOOS == "windows" {
			executable = "ffmpeg.exe"
		}
		destPath := filepath.Join(binDir, executable)
		
		if err := downloadFile(downloadURL, destPath, progressFn); err != nil {
			return fmt.Errorf("failed to download ffmpeg: %w", err)
		}
		
		if runtime.GOOS != "windows" {
			if err := os.Chmod(destPath, 0755); err != nil {
				return fmt.Errorf("failed to make ffmpeg executable: %w", err)
			}
		}
	}
	
	// Verify installation
	executable := "ffmpeg"
	if runtime.GOOS == "windows" {
		executable = "ffmpeg.exe"
	}
	destPath := filepath.Join(binDir, executable)
	
	if _, err := os.Stat(destPath); err != nil {
		return fmt.Errorf("ffmpeg installation verification failed: %w", err)
	}
	
	if progressFn != nil {
		progressFn(fmt.Sprintf("✓ ffmpeg installed at: %s", destPath))
	}
	
	return nil
}

// downloadFile downloads a file from url to filepath
func downloadFile(url, filepath string, progressFn func(string)) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	
	// Download with progress
	buf := make([]byte, 32*1024)
	var downloaded int64
	total := resp.ContentLength
	
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)
			
			if progressFn != nil && total > 0 {
				percentage := float64(downloaded) / float64(total) * 100
				progressFn(fmt.Sprintf("Downloading... %.1f%% (%d/%d MB)", 
					percentage, downloaded/(1024*1024), total/(1024*1024)))
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	
	return nil
}

// extractFFMPEGFromZip extracts ffmpeg binary from zip archive
func extractFFMPEGFromZip(zipPath, destDir string, progressFn func(string)) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	
	executable := "ffmpeg"
	if runtime.GOOS == "windows" {
		executable = "ffmpeg.exe"
	}
	
	// Find and extract ffmpeg binary
	for _, f := range r.File {
		// Look for ffmpeg binary in the archive
		if strings.Contains(f.Name, executable) && !strings.Contains(f.Name, "doc") {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			
			destPath := filepath.Join(destDir, executable)
			outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			
			if _, err := io.Copy(outFile, rc); err != nil {
				return err
			}
			
			if runtime.GOOS != "windows" {
				if err := os.Chmod(destPath, 0755); err != nil {
					return err
				}
			}
			
			return nil
		}
	}
	
	return fmt.Errorf("ffmpeg binary not found in archive")
}

// extractFFMPEGFromTar extracts ffmpeg binary from tar.gz or tar.xz archive
func extractFFMPEGFromTar(tarPath, destDir string, progressFn func(string)) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Handle gzip compression
	var reader io.Reader = file
	if strings.HasSuffix(tarPath, ".gz") {
		gzr, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzr.Close()
		reader = gzr
	}
	// Note: tar.xz would need xz decompression library
	// For production, consider using exec to call tar command
	
	tr := tar.NewReader(reader)
	
	executable := "ffmpeg"
	destPath := filepath.Join(destDir, executable)
	
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		
		// Look for ffmpeg binary
		if strings.Contains(header.Name, "ffmpeg") && 
		   !strings.Contains(header.Name, "doc") && 
		   header.Typeflag == tar.TypeReg {
			
			outFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer outFile.Close()
			
			if _, err := io.Copy(outFile, tr); err != nil {
				return err
			}
			
			if err := os.Chmod(destPath, 0755); err != nil {
				return err
			}
			
			return nil
		}
	}
	
	return fmt.Errorf("ffmpeg binary not found in archive")
}

// CheckInstallation verifies if yt-dlp and ffmpeg are installed
func CheckInstallation() (ytdlpInstalled, ffmpegInstalled bool, err error) {
	ytdlpPath, ytdlpErr := GetYTDLPPath()
	ffmpegPath, ffmpegErr := GetFFMPEGPath()
	
	ytdlpInstalled = ytdlpErr == nil && ytdlpPath != ""
	ffmpegInstalled = ffmpegErr == nil && ffmpegPath != ""
	
	return ytdlpInstalled, ffmpegInstalled, nil
}
