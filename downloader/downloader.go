package downloader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	
	"github.com/debargha2001/gostreampuller/internal/installer"
)

// Auto-installation state
var (
	installAttempted  bool
	installMutex      sync.Mutex
	autoInstallOnce   sync.Once
)

// tryGetLocalBinary attempts to find a locally installed binary
func tryGetLocalBinary(name string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return name // Fall back to system PATH
	}
	
	binDir := filepath.Join(homeDir, ".gostreampuller", "bin")
	
	// Check for executable with .exe on Windows
	var candidates []string
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		candidates = []string{
			filepath.Join(binDir, name+".exe"),
			filepath.Join(binDir, name),
		}
	} else {
		candidates = []string{
			filepath.Join(binDir, name),
		}
	}
	
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return name // Fall back to system PATH
}

func init() {
	// Auto-detect locally installed binaries on package initialization
	YTDLPPath = tryGetLocalBinary("yt-dlp")
	FFMPEGPath = tryGetLocalBinary("ffmpeg")
}

// checkBinaryExists verifies if a binary is executable
func checkBinaryExists(path string) bool {
	// If path is just a name (no directory), check if it's in PATH
	if !filepath.IsAbs(path) && !strings.Contains(path, string(filepath.Separator)) {
		_, err := exec.LookPath(path)
		return err == nil
	}
	
	// Check if file exists and is executable
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	
	// Check if it's a file (not directory)
	if info.IsDir() {
		return false
	}
	
	// On Unix systems, check if executable bit is set
	if info.Mode()&0111 != 0 {
		return true
	}
	
	// On Windows, .exe files are executable
	if strings.HasSuffix(strings.ToLower(path), ".exe") {
		return true
	}
	
	return false
}

// wasInstalledViaCLI checks if the user previously ran gostreampuller-cli setup
// This is determined by the presence of a marker file created by the CLI
func wasInstalledViaCLI() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	
	markerFile := filepath.Join(homeDir, ".gostreampuller", ".cli_installed")
	_, err = os.Stat(markerFile)
	return err == nil
}

// ensureBinariesInstalled checks if required binaries exist and installs them if needed
// This is called automatically on first use unless:
// - GOSTREAMPULLER_NO_AUTO_INSTALL=1 is set
// - Binaries were already installed via gostreampuller-cli setup
func ensureBinariesInstalled() error {
	// Only attempt installation once
	installMutex.Lock()
	if installAttempted {
		installMutex.Unlock()
		return nil
	}
	installAttempted = true
	installMutex.Unlock()
	
	// Check if auto-installation is disabled
	if os.Getenv("GOSTREAMPULLER_NO_AUTO_INSTALL") == "1" {
		return nil
	}
	
	// Check if binaries already exist and are executable
	ytdlpExists := checkBinaryExists(YTDLPPath)
	ffmpegExists := checkBinaryExists(FFMPEGPath)
	
	// If both exist, no installation needed
	if ytdlpExists && ffmpegExists {
		return nil
	}
	
	// Check if user already ran CLI setup but binaries are in PATH (not local)
	// If binaries are found in system PATH, don't auto-install
	if checkBinaryExists("yt-dlp") && checkBinaryExists("ffmpeg") {
		// User has system binaries, respect that choice
		return nil
	}
	
	// Check if CLI setup was already run (marker file exists)
	if wasInstalledViaCLI() {
		// User ran CLI setup but binaries are missing/broken
		// Don't auto-install, show helpful error message
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "[gostreampuller] ⚠ Warning: Binaries not found")
		fmt.Fprintln(os.Stderr, "[gostreampuller] It looks like you ran 'gostreampuller-cli setup' before,")
		fmt.Fprintln(os.Stderr, "[gostreampuller] but the binaries are missing or corrupted.")
		fmt.Fprintln(os.Stderr, "[gostreampuller] Please run: gostreampuller-cli setup")
		fmt.Fprintln(os.Stderr, "")
		return nil
	}
	
	// No CLI setup was done, and binaries are missing - auto-install
	return autoInstallBinaries(ytdlpExists, ffmpegExists)
}

// autoInstallBinaries performs the actual installation
func autoInstallBinaries(ytdlpExists, ffmpegExists bool) error {
	// Import the installer package
	// Note: This is imported here to avoid init-time side effects
	// The installer package is only loaded when needed
	
	// We need to dynamically import and use the installer
	// For now, we'll create a simple inline installer to avoid circular deps
	
	if !ytdlpExists || !ffmpegExists {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintln(os.Stderr, "  gostreampuller: First-time setup")
		fmt.Fprintln(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  Required binaries are being installed automatically.")
		fmt.Fprintln(os.Stderr, "  This is a one-time process and takes 1-3 minutes.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  To disable auto-installation, set:")
		fmt.Fprintln(os.Stderr, "    export GOSTREAMPULLER_NO_AUTO_INSTALL=1")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintln(os.Stderr, "")
	}
	
	// Try to install missing binaries
	if !ytdlpExists {
		fmt.Fprintln(os.Stderr, "[gostreampuller] Installing yt-dlp...")
		if err := installYTDLPAuto(); err != nil {
			fmt.Fprintf(os.Stderr, "[gostreampuller] ⚠ Warning: Could not auto-install yt-dlp: %v\n", err)
			fmt.Fprintln(os.Stderr, "[gostreampuller] Please install manually or run: gostreampuller-cli setup")
			fmt.Fprintln(os.Stderr, "[gostreampuller] Falling back to system yt-dlp (if available)")
		} else {
			fmt.Fprintln(os.Stderr, "[gostreampuller] ✓ yt-dlp installed successfully")
			// Update path to the newly installed binary
			YTDLPPath = tryGetLocalBinary("yt-dlp")
		}
	}
	
	if !ffmpegExists {
		fmt.Fprintln(os.Stderr, "[gostreampuller] Installing ffmpeg...")
		if err := installFFMPEGAuto(); err != nil {
			fmt.Fprintf(os.Stderr, "[gostreampuller] ⚠ Warning: Could not auto-install ffmpeg: %v\n", err)
			fmt.Fprintln(os.Stderr, "[gostreampuller] Please install manually or run: gostreampuller-cli setup")
			fmt.Fprintln(os.Stderr, "[gostreampuller] Falling back to system ffmpeg (if available)")
		} else {
			fmt.Fprintln(os.Stderr, "[gostreampuller] ✓ ffmpeg installed successfully")
			// Update path to the newly installed binary
			FFMPEGPath = tryGetLocalBinary("ffmpeg")
		}
	}
	
	if (!ytdlpExists || !ffmpegExists) {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "[gostreampuller] ✓ Setup complete! Subsequent calls will be fast.")
		fmt.Fprintln(os.Stderr, "")
		
		// Create marker file to indicate auto-installation was done
		// This prevents confusion if user later runs CLI setup
		createAutoInstallMarker()
	}
	
	return nil
}

// createAutoInstallMarker creates a marker file indicating auto-installation happened
func createAutoInstallMarker() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	
	configDir := filepath.Join(homeDir, ".gostreampuller")
	markerFile := filepath.Join(configDir, ".auto_installed")
	
	// Create empty marker file
	os.WriteFile(markerFile, []byte("auto-installed"), 0644)
}

// installYTDLPAuto installs yt-dlp automatically (inline implementation to avoid circular deps)
func installYTDLPAuto() error {
	progressFn := func(msg string) {
		if os.Getenv("GOSTREAMPULLER_VERBOSE") == "1" {
			fmt.Fprintf(os.Stderr, "[gostreampuller]   %s\n", msg)
		}
	}
	
	return installer.InstallYTDLP(progressFn)
}

// installFFMPEGAuto installs ffmpeg automatically
func installFFMPEGAuto() error {
	progressFn := func(msg string) {
		if os.Getenv("GOSTREAMPULLER_VERBOSE") == "1" {
			fmt.Fprintf(os.Stderr, "[gostreampuller]   %s\n", msg)
		}
	}
	
	return installer.InstallFFMPEG(progressFn)
}


var (
	// YTDLPPath is the path to yt-dlp binary
	// Automatically set to local binary (~/.gostreampuller/bin/yt-dlp) if found,
	// otherwise falls back to system PATH ("yt-dlp")
	// Can be overridden using SetYTDLPPath()
	YTDLPPath string
	
	// FFMPEGPath is the path to ffmpeg binary
	// Automatically set to local binary (~/.gostreampuller/bin/ffmpeg) if found,
	// otherwise falls back to system PATH ("ffmpeg")
	// Can be overridden using SetFFMPEGPath()
	FFMPEGPath string
	
	// ChunkSize defines the buffer size for streaming operations (default: 32MB)
	ChunkSize = 32 * 1024 * 1024
	
	// MaxConcurrentDownloads limits parallel download operations (default: 3)
	MaxConcurrentDownloads = 3
)

// SetYTDLPPath sets a custom path for the yt-dlp binary.
// This overrides the auto-detected local binary.
// Use this if you want to use a specific yt-dlp installation.
//
// Example:
//   downloader.SetYTDLPPath("/usr/local/bin/yt-dlp")
//   downloader.SetYTDLPPath("C:\\tools\\yt-dlp.exe")  // Windows
func SetYTDLPPath(path string) {
	YTDLPPath = path
}

// SetFFMPEGPath sets a custom path for the ffmpeg binary.
// This overrides the auto-detected local binary.
// Use this if you want to use a specific ffmpeg installation.
//
// Example:
//   downloader.SetFFMPEGPath("/usr/local/bin/ffmpeg")
//   downloader.SetFFMPEGPath("C:\\ffmpeg\\bin\\ffmpeg.exe")  // Windows
func SetFFMPEGPath(path string) {
	FFMPEGPath = path
}

// ResetBinaryPaths resets binary paths to auto-detected defaults.
// Call this to revert any custom paths set by SetYTDLPPath() or SetFFMPEGPath().
func ResetBinaryPaths() {
	YTDLPPath = tryGetLocalBinary("yt-dlp")
	FFMPEGPath = tryGetLocalBinary("ffmpeg")
}

// SetChunkSize sets the buffer size for streaming operations
func SetChunkSize(size int) {
	if size > 0 {
		ChunkSize = size
	}
}

// SetMaxConcurrentDownloads sets the maximum number of concurrent downloads
func SetMaxConcurrentDownloads(max int) {
	if max > 0 {
		MaxConcurrentDownloads = max
	}
}

// DownloadProgress represents download progress information
type DownloadProgress struct {
	BytesDownloaded int64
	TotalBytes      int64
	Percentage      float64
	Stage           string
}

// ProgressCallback is called during download to report progress
type ProgressCallback func(progress DownloadProgress)

// VideoMetadata represents comprehensive metadata for a video
type VideoMetadata struct {
	// Basic Information
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Uploader    string `json:"uploader"`
	UploaderID  string `json:"uploader_id"`
	UploaderURL string `json:"uploader_url"`
	
	// Video Details
	Duration       int     `json:"duration"`        // Duration in seconds
	DurationString string  `json:"duration_string"` // Human readable duration
	ViewCount      int64   `json:"view_count"`
	LikeCount      int64   `json:"like_count"`
	DislikeCount   int64   `json:"dislike_count"`
	AverageRating  float64 `json:"average_rating"`
	
	// Technical Details
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	FPS            float64 `json:"fps"`
	VideoCodec     string  `json:"vcodec"`
	AudioCodec     string  `json:"acodec"`
	Format         string  `json:"format"`
	FormatID       string  `json:"format_id"`
	Extension      string  `json:"ext"`
	Filesize       int64   `json:"filesize"`
	FilesizeApprox int64   `json:"filesize_approx"`
	
	// URLs
	URL         string `json:"url"`
	WebpageURL  string `json:"webpage_url"`
	Thumbnail   string `json:"thumbnail"`
	
	// Timestamps
	UploadDate string `json:"upload_date"`
	ReleaseDate string `json:"release_date"`
	Timestamp  int64  `json:"timestamp"`
	
	// Additional Info
	Categories  []string               `json:"categories"`
	Tags        []string               `json:"tags"`
	IsLive      bool                   `json:"is_live"`
	WasLive     bool                   `json:"was_live"`
	LiveStatus  string                 `json:"live_status"`
	Channel     string                 `json:"channel"`
	ChannelID   string                 `json:"channel_id"`
	ChannelURL  string                 `json:"channel_url"`
	Subtitles   map[string]interface{} `json:"subtitles"`
	
	// Platform Specific
	Extractor    string                 `json:"extractor"`
	ExtractorKey string                 `json:"extractor_key"`
	
	// Raw metadata for additional fields
	Raw map[string]interface{} `json:"-"`
}

// GetVideoMetadata fetches comprehensive metadata for a video without downloading it
// Returns detailed information about the video including title, duration, formats, quality, etc.
//
// Example:
//   metadata, err := downloader.GetVideoMetadata("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Title: %s\nDuration: %s\nViews: %d\n", 
//       metadata.Title, metadata.DurationString, metadata.ViewCount)
func GetVideoMetadata(url string) (*VideoMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	
	return GetVideoMetadataWithContext(ctx, url)
}

// GetVideoMetadataWithContext fetches video metadata with a custom context for timeout/cancellation
func GetVideoMetadataWithContext(ctx context.Context, url string) (*VideoMetadata, error) {
	// Auto-install binaries if needed (only happens once)
	if err := ensureBinariesInstalled(); err != nil {
		return nil, fmt.Errorf("failed to ensure binaries are installed: %w", err)
	}
	
	// Use yt-dlp with --dump-json to get metadata without downloading
	cmd := exec.CommandContext(ctx, YTDLPPath, 
		"--dump-json",
		"--no-playlist",
		"--no-warnings",
		url,
	)
	
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("failed to fetch metadata: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute yt-dlp: %w", err)
	}
	
	// Parse JSON output
	var rawMetadata map[string]interface{}
	if err := json.Unmarshal(output, &rawMetadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}
	
	// Create VideoMetadata struct
	metadata := &VideoMetadata{
		Raw: rawMetadata,
	}
	
	// Marshal and unmarshal to populate struct fields
	if err := json.Unmarshal(output, metadata); err != nil {
		return nil, fmt.Errorf("failed to map metadata: %w", err)
	}
	
	return metadata, nil
}

// streamCommand executes a command and streams its output to handle large files
func streamCommand(ctx context.Context, cmd *exec.Cmd, progressCb ProgressCallback, stage string) error {
	var wg sync.WaitGroup
	var errOut error
	var mu sync.Mutex

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Stream stdout in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, ChunkSize), ChunkSize)
		
		for scanner.Scan() {
			// Parse progress from output if callback provided
			if progressCb != nil {
				line := scanner.Text()
				// yt-dlp outputs progress information that can be parsed
				if strings.Contains(line, "%") || strings.Contains(line, "ETA") {
					progressCb(DownloadProgress{
						Stage: stage,
					})
				}
			}
		}
		
		if err := scanner.Err(); err != nil && err != io.EOF {
			mu.Lock()
			if errOut == nil {
				errOut = fmt.Errorf("stdout scan error: %w", err)
			}
			mu.Unlock()
		}
	}()

	// Stream stderr in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		scanner.Buffer(make([]byte, ChunkSize), ChunkSize)
		
		for scanner.Scan() {
			// Log errors but don't fail on warnings
			_ = scanner.Text()
		}
		
		if err := scanner.Err(); err != nil && err != io.EOF {
			mu.Lock()
			if errOut == nil {
				errOut = fmt.Errorf("stderr scan error: %w", err)
			}
			mu.Unlock()
		}
	}()

	// Wait for streams to complete
	wg.Wait()

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		if errOut != nil {
			return fmt.Errorf("command failed: %v, %w", err, errOut)
		}
		return fmt.Errorf("command failed: %w", err)
	}

	return errOut
}

// copyFileStreaming copies a file using streaming to handle large files efficiently
func copyFileStreaming(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Use buffered I/O for better performance with large files
	buf := make([]byte, ChunkSize)
	written, err := io.CopyBuffer(destFile, sourceFile, buf)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	_ = written
	return nil
}

// DownloadVideo downloads a video, allowing optional format, resolution, and codec parameters.
// If any parameter is empty, defaults will be used.
// This function uses streaming and concurrent processing to handle large files efficiently.
// Files are saved to the current working directory.
func DownloadVideo(url string, format string, resolution string, codec string) (string, error) {
	return DownloadVideoWithProgress(url, format, resolution, codec, nil)
}

// DownloadVideoToDir downloads a video to a specific directory.
// If outputDir is empty, files are saved to the current working directory.
func DownloadVideoToDir(url string, format string, resolution string, codec string, outputDir string) (string, error) {
	return DownloadVideoToDirWithProgress(url, format, resolution, codec, outputDir, nil)
}

// DownloadVideoWithProgress downloads a video with progress callback support.
// The progressCb function is called periodically with download progress information.
// Files are saved to the current working directory.
func DownloadVideoWithProgress(url string, format string, resolution string, codec string, progressCb ProgressCallback) (string, error) {
	return DownloadVideoToDirWithProgress(url, format, resolution, codec, "", progressCb)
}

// DownloadVideoToDirWithProgress downloads a video to a specific directory with progress callback support.
// If outputDir is empty, files are saved to the current working directory.
func DownloadVideoToDirWithProgress(url string, format string, resolution string, codec string, outputDir string, progressCb ProgressCallback) (string, error) {
	// Auto-install binaries if needed (only happens once)
	if err := ensureBinariesInstalled(); err != nil {
		return "", fmt.Errorf("failed to ensure binaries are installed: %w", err)
	}
	
	if format == "" {
		format = "mp4"
	}
	if resolution == "" {
		resolution = "720"
	}
	if codec == "" {
		codec = "avc1"
	}
	
	// Use custom output directory if provided
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	filename := fmt.Sprintf("video_%d.%%(ext)s", time.Now().UnixNano())
	var temp string
	if outputDir != "" {
		temp = filepath.Join(outputDir, filename)
	} else {
		temp = filename
	}
	selector := fmt.Sprintf("bestvideo[height<=%s][vcodec*=%s]+bestaudio/best", resolution, codec)

	// Use yt-dlp with options optimized for large files
	cmd := exec.CommandContext(ctx, YTDLPPath,
		"-f", selector,
		"-o", temp,
		"--no-part",                    // Don't use .part files for large downloads
		"--concurrent-fragments", "3",   // Download fragments concurrently
		"--buffer-size", "32K",          // Set buffer size
		"--retries", "10",               // Retry on failure
		"--fragment-retries", "10",      // Retry fragments
		url,
	)

	if progressCb != nil {
		progressCb(DownloadProgress{Stage: "Downloading video"})
	}

	if err := streamCommand(ctx, cmd, progressCb, "downloading"); err != nil {
		return "", fmt.Errorf("yt-dlp video download failed: %w", err)
	}

	// Find the actual downloaded file by checking common extensions
	var downloaded string
	possibleExtensions := []string{"mkv", "mp4", "webm", "avi", "mov", "flv"}

	for _, ext := range possibleExtensions {
		candidate := strings.Replace(temp, "%(ext)s", ext, 1)
		if _, err := os.Stat(candidate); err == nil {
			downloaded = candidate
			break
		}
	}

	if downloaded == "" {
		return "", fmt.Errorf("could not find downloaded video file")
	}

	// If format is different from downloaded format, convert it
	finalOutput := strings.Replace(temp, "%(ext)s", format, 1)
	if downloaded != finalOutput {
		if progressCb != nil {
			progressCb(DownloadProgress{Stage: "Converting video format"})
		}

		convertCtx, convertCancel := context.WithTimeout(context.Background(), 20*time.Minute)
		defer convertCancel()

		// Use streaming copy for format conversion to handle large files
		ffmpeg := exec.CommandContext(convertCtx, FFMPEGPath,
			"-i", downloaded,
			"-c", "copy",
			"-movflags", "+faststart",     // Optimize for streaming
			"-max_muxing_queue_size", "1024", // Handle large files
			"-y",
			finalOutput,
		)

		if err := streamCommand(convertCtx, ffmpeg, progressCb, "converting"); err != nil {
			return "", fmt.Errorf("ffmpeg conversion failed: %w", err)
		}
		defer os.Remove(downloaded)
		
		if progressCb != nil {
			progressCb(DownloadProgress{Stage: "Completed", Percentage: 100.0})
		}
		
		return filepath.Abs(finalOutput)
	}

	if progressCb != nil {
		progressCb(DownloadProgress{Stage: "Completed", Percentage: 100.0})
	}

	return filepath.Abs(downloaded)
}

// DownloadAudio downloads audio, allowing optional output format, codec, and bitrate parameters.
// If any parameter is empty, defaults will be used.
// This function uses streaming and concurrent processing to handle large files efficiently.
// Files are saved to the current working directory.
func DownloadAudio(url string, outputFormat string, codec string, bitrate string) (string, error) {
	return DownloadAudioWithProgress(url, outputFormat, codec, bitrate, nil)
}

// DownloadAudioToDir downloads audio to a specific directory.
// If outputDir is empty, files are saved to the current working directory.
func DownloadAudioToDir(url string, outputFormat string, codec string, bitrate string, outputDir string) (string, error) {
	return DownloadAudioToDirWithProgress(url, outputFormat, codec, bitrate, outputDir, nil)
}

// DownloadAudioWithProgress downloads audio with progress callback support.
// The progressCb function is called periodically with download progress information.
// Files are saved to the current working directory.
func DownloadAudioWithProgress(url string, outputFormat string, codec string, bitrate string, progressCb ProgressCallback) (string, error) {
	return DownloadAudioToDirWithProgress(url, outputFormat, codec, bitrate, "", progressCb)
}

// DownloadAudioToDirWithProgress downloads audio to a specific directory with progress callback support.
// If outputDir is empty, files are saved to the current working directory.
func DownloadAudioToDirWithProgress(url string, outputFormat string, codec string, bitrate string, outputDir string, progressCb ProgressCallback) (string, error) {
	// Auto-install binaries if needed (only happens once)
	if err := ensureBinariesInstalled(); err != nil {
		return "", fmt.Errorf("failed to ensure binaries are installed: %w", err)
	}
	
	if outputFormat == "" {
		outputFormat = "mp3"
	}
	if codec == "" {
		codec = "libmp3lame"
	}
	if bitrate == "" {
		bitrate = "128k"
	}
	
	// Use custom output directory if provided
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	filename := fmt.Sprintf("audio_%d.%%(ext)s", time.Now().UnixNano())
	var temp string
	if outputDir != "" {
		temp = filepath.Join(outputDir, filename)
	} else {
		temp = filename
	}
	
	// Use yt-dlp with options optimized for large files
	cmd := exec.CommandContext(ctx, YTDLPPath,
		"-f", "bestaudio",
		"-o", temp,
		"--no-part",                    // Don't use .part files
		"--concurrent-fragments", "3",   // Download fragments concurrently
		"--buffer-size", "32K",          // Set buffer size
		"--retries", "10",               // Retry on failure
		"--fragment-retries", "10",      // Retry fragments
		url,
	)

	if progressCb != nil {
		progressCb(DownloadProgress{Stage: "Downloading audio"})
	}

	if err := streamCommand(ctx, cmd, progressCb, "downloading"); err != nil {
		return "", fmt.Errorf("yt-dlp audio fetch failed: %w", err)
	}

	// Find the downloaded file (could be webm, m4a, opus, etc.)
	possibleExtensions := []string{"webm", "m4a", "opus", "ogg", "mp3", "aac"}
	var original string
	
	for _, ext := range possibleExtensions {
		candidate := strings.Replace(temp, "%(ext)s", ext, 1)
		if _, err := os.Stat(candidate); err == nil {
			original = candidate
			break
		}
	}

	if original == "" {
		return "", fmt.Errorf("could not find downloaded audio file")
	}

	output := strings.Replace(temp, "%(ext)s", outputFormat, 1)

	if progressCb != nil {
		progressCb(DownloadProgress{Stage: "Converting audio format"})
	}

	convertCtx, convertCancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer convertCancel()

	// Use streaming conversion for large audio files
	ffmpeg := exec.CommandContext(convertCtx, FFMPEGPath,
		"-i", original,
		"-vn",
		"-acodec", codec,
		"-ab", bitrate,
		"-max_muxing_queue_size", "1024", // Handle large files
		"-y",
		output,
	)

	if err := streamCommand(convertCtx, ffmpeg, progressCb, "converting"); err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	defer os.Remove(original)
	
	if progressCb != nil {
		progressCb(DownloadProgress{Stage: "Completed", Percentage: 100.0})
	}
	
	return filepath.Abs(output)
}
