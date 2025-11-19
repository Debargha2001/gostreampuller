package tests

import (
	"context"
	"testing"
	"time"

	"github.com/debargha2001/gostreampuller/downloader"
)

func TestGetVideoMetadata(t *testing.T) {
	// Use a short, stable video for testing
	testURL := "https://www.youtube.com/watch?v=jNQXAC9IVRw" // "Me at the zoo" - first YouTube video
	
	metadata, err := downloader.GetVideoMetadata(testURL)
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}
	
	// Verify essential fields are populated
	if metadata.ID == "" {
		t.Error("Expected video ID to be populated")
	}
	
	if metadata.Title == "" {
		t.Error("Expected title to be populated")
	}
	
	if metadata.Uploader == "" {
		t.Error("Expected uploader to be populated")
	}
	
	if metadata.Duration <= 0 {
		t.Error("Expected duration to be greater than 0")
	}
	
	if metadata.WebpageURL == "" {
		t.Error("Expected webpage URL to be populated")
	}
	
	t.Logf("Video Metadata:")
	t.Logf("  ID: %s", metadata.ID)
	t.Logf("  Title: %s", metadata.Title)
	t.Logf("  Uploader: %s", metadata.Uploader)
	t.Logf("  Duration: %s (%d seconds)", metadata.DurationString, metadata.Duration)
	t.Logf("  Views: %d", metadata.ViewCount)
	t.Logf("  Format: %s", metadata.Format)
	t.Logf("  Resolution: %dx%d", metadata.Width, metadata.Height)
	t.Logf("  Video Codec: %s", metadata.VideoCodec)
	t.Logf("  Audio Codec: %s", metadata.AudioCodec)
}

func TestGetVideoMetadataWithContext(t *testing.T) {
	testURL := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	// Test with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	metadata, err := downloader.GetVideoMetadataWithContext(ctx, testURL)
	if err != nil {
		t.Fatalf("Failed to get metadata with context: %v", err)
	}
	
	if metadata.Title == "" {
		t.Error("Expected title to be populated")
	}
}

func TestGetVideoMetadataInvalidURL(t *testing.T) {
	_, err := downloader.GetVideoMetadata("https://invalid-url-that-does-not-exist.com/video")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
	t.Logf("Got expected error: %v", err)
}

func TestGetVideoMetadataTimeout(t *testing.T) {
	testURL := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	// Test with very short timeout (should fail)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	time.Sleep(10 * time.Millisecond) // Ensure context expires
	
	_, err := downloader.GetVideoMetadataWithContext(ctx, testURL)
	if err == nil {
		t.Error("Expected timeout error")
	}
	t.Logf("Got expected timeout error: %v", err)
}

func TestGetVideoMetadataMultiplePlatforms(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		platform string
	}{
		{
			name:     "YouTube",
			url:      "https://www.youtube.com/watch?v=jNQXAC9IVRw",
			platform: "youtube",
		},
		// Add more platforms as needed for testing
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := downloader.GetVideoMetadata(tt.url)
			if err != nil {
				t.Fatalf("Failed to get %s metadata: %v", tt.platform, err)
			}
			
			if metadata.ExtractorKey == "" {
				t.Error("Expected extractor key to be populated")
			}
			
			t.Logf("%s Video:", tt.platform)
			t.Logf("  Title: %s", metadata.Title)
			t.Logf("  Extractor: %s", metadata.ExtractorKey)
			t.Logf("  Duration: %s", metadata.DurationString)
		})
	}
}

func TestGetVideoMetadataLiveStream(t *testing.T) {
	t.Skip("Skipping live stream test - requires active live stream URL")
	
	// Example test for live streams (uncomment and provide URL when testing)
	// liveURL := "https://www.youtube.com/watch?v=LIVE_STREAM_ID"
	// metadata, err := downloader.GetVideoMetadata(liveURL)
	// if err != nil {
	//     t.Fatalf("Failed to get live stream metadata: %v", err)
	// }
	// 
	// if !metadata.IsLive {
	//     t.Error("Expected IsLive to be true for live stream")
	// }
}

func TestVideoMetadataFields(t *testing.T) {
	testURL := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	metadata, err := downloader.GetVideoMetadata(testURL)
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}
	
	// Test that Raw field contains additional data
	if metadata.Raw == nil {
		t.Error("Expected Raw metadata map to be populated")
	}
	
	// Verify some fields that should always be present
	tests := []struct {
		field string
		value interface{}
	}{
		{"ID", metadata.ID},
		{"Title", metadata.Title},
		{"WebpageURL", metadata.WebpageURL},
	}
	
	for _, tt := range tests {
		if tt.value == nil || tt.value == "" {
			t.Errorf("Field %s should not be empty", tt.field)
		}
	}
}
