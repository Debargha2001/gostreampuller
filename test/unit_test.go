package tests

import (
	"testing"

	"github.com/debargha2001/gostreampuller/downloader"
)

// TestPackageInitialization tests package initialization
func TestPackageInitialization(t *testing.T) {
	if downloader.ChunkSize != 32*1024*1024 {
		t.Errorf("Expected default chunk size 32MB, got %d", downloader.ChunkSize)
	}
	
	if downloader.MaxConcurrentDownloads != 3 {
		t.Errorf("Expected default max concurrent downloads 3, got %d", downloader.MaxConcurrentDownloads)
	}
	
	if downloader.YTDLPPath == "" {
		t.Error("YTDLPPath not initialized")
	}
	
	if downloader.FFMPEGPath == "" {
		t.Error("FFMPEGPath not initialized")
	}
	
	t.Log("✓ Package initialization successful")
}

// TestSetters tests all setter functions
func TestSetters(t *testing.T) {
	// Save original values
	origChunkSize := downloader.ChunkSize
	origMaxConcurrent := downloader.MaxConcurrentDownloads
	origYTDLP := downloader.YTDLPPath
	origFFMPEG := downloader.FFMPEGPath
	
	// Restore at end
	defer func() {
		downloader.ChunkSize = origChunkSize
		downloader.MaxConcurrentDownloads = origMaxConcurrent
		downloader.YTDLPPath = origYTDLP
		downloader.FFMPEGPath = origFFMPEG
	}()
	
	// Test SetChunkSize
	downloader.SetChunkSize(64 * 1024 * 1024)
	if downloader.ChunkSize != 64*1024*1024 {
		t.Error("SetChunkSize failed")
	}
	
	// Test SetMaxConcurrentDownloads
	downloader.SetMaxConcurrentDownloads(10)
	if downloader.MaxConcurrentDownloads != 10 {
		t.Error("SetMaxConcurrentDownloads failed")
	}
	
	// Test SetYTDLPPath (custom path should override auto-detected)
	customYTDLP := "/custom/path/yt-dlp"
	downloader.SetYTDLPPath(customYTDLP)
	if downloader.YTDLPPath != customYTDLP {
		t.Error("SetYTDLPPath failed")
	}
	
	// Test SetFFMPEGPath (custom path should override auto-detected)
	customFFMPEG := "/custom/path/ffmpeg"
	downloader.SetFFMPEGPath(customFFMPEG)
	if downloader.FFMPEGPath != customFFMPEG {
		t.Error("SetFFMPEGPath failed")
	}
	
	// Test ResetBinaryPaths (should revert to auto-detected paths)
	downloader.ResetBinaryPaths()
	if downloader.YTDLPPath == customYTDLP {
		t.Error("ResetBinaryPaths failed to reset YTDLPPath")
	}
	if downloader.FFMPEGPath == customFFMPEG {
		t.Error("ResetBinaryPaths failed to reset FFMPEGPath")
	}
	
	// Verify reset paths are not empty (should have auto-detected values)
	if downloader.YTDLPPath == "" {
		t.Error("YTDLPPath empty after reset")
	}
	if downloader.FFMPEGPath == "" {
		t.Error("FFMPEGPath empty after reset")
	}
	
	t.Log("✓ All setter functions work correctly")
	t.Logf("  Custom paths can override auto-detected defaults")
	t.Logf("  ResetBinaryPaths() restores auto-detection")
}

// TestDownloadProgressType tests the DownloadProgress type
func TestDownloadProgressType(t *testing.T) {
	progress := downloader.DownloadProgress{
		BytesDownloaded: 1024 * 1024,
		TotalBytes:      10 * 1024 * 1024,
		Percentage:      10.0,
		Stage:           "Downloading",
	}
	
	if progress.BytesDownloaded != 1024*1024 {
		t.Error("BytesDownloaded field not working")
	}
	
	if progress.TotalBytes != 10*1024*1024 {
		t.Error("TotalBytes field not working")
	}
	
	if progress.Percentage != 10.0 {
		t.Error("Percentage field not working")
	}
	
	if progress.Stage != "Downloading" {
		t.Error("Stage field not working")
	}
	
	t.Log("✓ DownloadProgress type works correctly")
}
