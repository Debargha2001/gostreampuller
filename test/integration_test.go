package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/debargha2001/gostreampuller/downloader"
)

// TestDownloadVideo tests basic video download functionality
func TestDownloadVideo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This is a short test video - "Me at the zoo" (first YouTube video)
	url := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	video, err := downloader.DownloadVideo(url, "mp4", "360", "avc1")
	if err != nil {
		t.Fatalf("Failed to download video: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(video); os.IsNotExist(err) {
		t.Fatalf("Downloaded video file does not exist: %s", video)
	}

	// Check file size (should be > 0)
	fileInfo, err := os.Stat(video)
	if err != nil {
		t.Fatalf("Failed to stat video file: %v", err)
	}
	
	if fileInfo.Size() == 0 {
		t.Error("Downloaded video file is empty")
	}

	// Clean up
	defer os.Remove(video)

	t.Logf("✓ Successfully downloaded video: %s (%.2f MB)", 
		video, float64(fileInfo.Size())/(1024*1024))
}

// TestDownloadVideoWithProgress tests video download with progress callback
func TestDownloadVideoWithProgress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	progressCalled := false
	completedCalled := false
	stages := []string{}
	
	callback := func(progress downloader.DownloadProgress) {
		progressCalled = true
		stages = append(stages, progress.Stage)
		t.Logf("Progress: %s (%.1f%%)", progress.Stage, progress.Percentage)
		
		if progress.Stage == "Completed" {
			completedCalled = true
		}
	}

	video, err := downloader.DownloadVideoWithProgress(url, "mp4", "360", "avc1", callback)
	if err != nil {
		t.Fatalf("Failed to download video: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(video); os.IsNotExist(err) {
		t.Fatalf("Downloaded video file does not exist: %s", video)
	}

	// Verify callback was called
	if !progressCalled {
		t.Error("Progress callback was never called")
	}

	if !completedCalled {
		t.Error("Completion callback was never called")
	}
	
	// Verify we got different stages
	if len(stages) < 2 {
		t.Errorf("Expected multiple progress stages, got: %v", stages)
	}

	// Clean up
	defer os.Remove(video)

	t.Logf("✓ Successfully downloaded video with progress tracking")
}

// TestDownloadAudio tests basic audio download functionality
func TestDownloadAudio(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	audio, err := downloader.DownloadAudio(url, "mp3", "libmp3lame", "128k")
	if err != nil {
		t.Fatalf("Failed to download audio: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(audio); os.IsNotExist(err) {
		t.Fatalf("Downloaded audio file does not exist: %s", audio)
	}

	// Verify it's an MP3 file
	if !strings.HasSuffix(audio, ".mp3") {
		t.Errorf("Expected MP3 file, got: %s", audio)
	}
	
	// Check file size
	fileInfo, err := os.Stat(audio)
	if err != nil {
		t.Fatalf("Failed to stat audio file: %v", err)
	}
	
	if fileInfo.Size() == 0 {
		t.Error("Downloaded audio file is empty")
	}

	// Clean up
	defer os.Remove(audio)

	t.Logf("✓ Successfully downloaded audio: %s (%.2f MB)", 
		audio, float64(fileInfo.Size())/(1024*1024))
}

// TestDownloadAudioWithProgress tests audio download with progress callback
func TestDownloadAudioWithProgress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	stages := []string{}
	
	callback := func(progress downloader.DownloadProgress) {
		stages = append(stages, progress.Stage)
		t.Logf("Progress: %s", progress.Stage)
	}

	audio, err := downloader.DownloadAudioWithProgress(url, "mp3", "libmp3lame", "192k", callback)
	if err != nil {
		t.Fatalf("Failed to download audio: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(audio); os.IsNotExist(err) {
		t.Fatalf("Downloaded audio file does not exist: %s", audio)
	}

	// Verify callback was called with different stages
	if len(stages) == 0 {
		t.Error("Progress callback was never called")
	}

	// Clean up
	defer os.Remove(audio)

	t.Logf("✓ Successfully downloaded audio with %d progress updates", len(stages))
}

// TestCustomChunkSize tests setting custom chunk size
func TestCustomChunkSize(t *testing.T) {
	originalSize := downloader.ChunkSize
	defer func() { downloader.ChunkSize = originalSize }()

	downloader.SetChunkSize(64 * 1024 * 1024)
	
	if downloader.ChunkSize != 64*1024*1024 {
		t.Errorf("Expected chunk size %d, got %d", 64*1024*1024, downloader.ChunkSize)
	}

	// Test invalid size
	downloader.SetChunkSize(-1)
	if downloader.ChunkSize != 64*1024*1024 {
		t.Error("Chunk size should not change with invalid input")
	}

	downloader.SetChunkSize(0)
	if downloader.ChunkSize != 64*1024*1024 {
		t.Error("Chunk size should not change with zero input")
	}
	
	t.Log("✓ Chunk size configuration works correctly")
}

// TestCustomConcurrency tests setting custom concurrent downloads
func TestCustomConcurrency(t *testing.T) {
	originalMax := downloader.MaxConcurrentDownloads
	defer func() { downloader.MaxConcurrentDownloads = originalMax }()

	downloader.SetMaxConcurrentDownloads(5)
	
	if downloader.MaxConcurrentDownloads != 5 {
		t.Errorf("Expected max concurrent downloads %d, got %d", 5, downloader.MaxConcurrentDownloads)
	}

	// Test invalid value
	downloader.SetMaxConcurrentDownloads(-1)
	if downloader.MaxConcurrentDownloads != 5 {
		t.Error("Max concurrent downloads should not change with invalid input")
	}

	downloader.SetMaxConcurrentDownloads(0)
	if downloader.MaxConcurrentDownloads != 5 {
		t.Error("Max concurrent downloads should not change with zero input")
	}
	
	t.Log("✓ Concurrency configuration works correctly")
}

// TestInvalidURL tests handling of invalid URLs
func TestInvalidURL(t *testing.T) {
	_, err := downloader.DownloadVideo("https://invalid-url-that-does-not-exist-12345.com/video", "mp4", "720", "avc1")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
	t.Logf("✓ Correctly handled invalid URL: %v", err)
}

// TestDefaultParameters tests that default parameters work correctly
func TestDefaultParameters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	// Test with empty parameters (should use defaults)
	video, err := downloader.DownloadVideo(url, "", "", "")
	if err != nil {
		t.Fatalf("Failed to download with defaults: %v", err)
	}

	defer os.Remove(video)

	// Should default to mp4
	if !strings.HasSuffix(video, ".mp4") {
		t.Errorf("Expected default format mp4, got: %s", video)
	}

	t.Log("✓ Default parameters work correctly")
}

// TestConcurrentDownloads tests multiple concurrent downloads
func TestConcurrentDownloads(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	urls := []string{
		"https://www.youtube.com/watch?v=jNQXAC9IVRw",
		"https://www.youtube.com/watch?v=jNQXAC9IVRw", // Same video, different download
	}

	results := make(chan string, len(urls))
	errors := make(chan error, len(urls))

	for i, url := range urls {
		go func(index int, videoURL string) {
			video, err := downloader.DownloadVideo(videoURL, "mp4", "360", "avc1")
			if err != nil {
				errors <- err
				return
			}
			results <- video
		}(i, url)
	}

	downloaded := []string{}
	
	// Collect results
	for i := 0; i < len(urls); i++ {
		select {
		case video := <-results:
			downloaded = append(downloaded, video)
			t.Logf("Downloaded: %s", video)
		case err := <-errors:
			t.Errorf("Error in concurrent download: %v", err)
		case <-time.After(5 * time.Minute):
			t.Fatal("Timeout waiting for concurrent downloads")
		}
	}

	// Clean up
	for _, video := range downloaded {
		os.Remove(video)
	}

	if len(downloaded) != len(urls) {
		t.Errorf("Expected %d downloads, got %d", len(urls), len(downloaded))
	}
	
	t.Logf("✓ Successfully completed %d concurrent downloads", len(downloaded))
}

// TestBinaryPaths tests that binary paths are correctly configured
func TestBinaryPaths(t *testing.T) {
	if downloader.YTDLPPath == "" {
		t.Error("YTDLPPath is empty")
	}
	
	if downloader.FFMPEGPath == "" {
		t.Error("FFMPEGPath is empty")
	}
	
	t.Logf("✓ Binary paths configured:")
	t.Logf("  yt-dlp: %s", downloader.YTDLPPath)
	t.Logf("  ffmpeg: %s", downloader.FFMPEGPath)
}

// Example demonstrates basic usage
func ExampleDownloadVideo() {
	video, err := downloader.DownloadVideo(
		"https://youtube.com/watch?v=example",
		"mp4",
		"720",
		"avc1",
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Downloaded: %s\n", video)
}

// Example demonstrates usage with progress callback
func ExampleDownloadVideoWithProgress() {
	callback := func(progress downloader.DownloadProgress) {
		fmt.Printf("Stage: %s (%.1f%%)\n", progress.Stage, progress.Percentage)
	}

	video, err := downloader.DownloadVideoWithProgress(
		"https://youtube.com/watch?v=example",
		"mp4",
		"1080",
		"avc1",
		callback,
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Downloaded: %s\n", video)
}
