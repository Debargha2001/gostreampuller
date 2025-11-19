package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/debargha2001/gostreampuller/downloader"
)

func main() {
	url := "https://www.youtube.com/watch?v=jNQXAC9IVRw" // Short test video

	// Example 1: Download to current directory (default behavior)
	fmt.Println("=== Example 1: Default - Current Directory ===")
	video1, err := downloader.DownloadVideo(url, "mp4", "360", "avc1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Downloaded to: %s\n\n", video1)

	// Example 2: Download to custom directory
	fmt.Println("=== Example 2: Custom Directory ===")
	customDir := filepath.Join(os.TempDir(), "gostreampuller_test")
	video2, err := downloader.DownloadVideoToDir(url, "mp4", "360", "avc1", customDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Downloaded to: %s\n", video2)
	fmt.Printf("✓ Directory exists: %s\n\n", customDir)

	// Example 3: Download audio to custom directory with progress tracking
	fmt.Println("=== Example 3: Audio with Progress ===")
	musicDir := filepath.Join(os.TempDir(), "gostreampuller_music")
	
	progressCallback := func(progress downloader.DownloadProgress) {
		if progress.Percentage > 0 {
			fmt.Printf("[%s] %.1f%%\n", progress.Stage, progress.Percentage)
		} else {
			fmt.Printf("[%s]\n", progress.Stage)
		}
	}

	audio, err := downloader.DownloadAudioToDirWithProgress(
		url, "mp3", "libmp3lame", "192k",
		musicDir,
		progressCallback,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Downloaded to: %s\n\n", audio)

	// Example 4: Using empty string for current directory (same as default)
	fmt.Println("=== Example 4: Empty String = Current Directory ===")
	video3, err := downloader.DownloadVideoToDir(url, "mp4", "360", "avc1", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Downloaded to: %s\n\n", video3)

	// Cleanup
	fmt.Println("=== Cleanup ===")
	os.Remove(video1)
	os.Remove(video3)
	os.RemoveAll(customDir)
	os.RemoveAll(musicDir)
	fmt.Println("✓ Cleaned up test files")
}
