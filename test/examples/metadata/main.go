package main

import (
	"fmt"
	"log"

	"github.com/debargha2001/gostreampuller/downloader"
)

func main() {
	// Example 1: Get basic video information
	fmt.Println("=== Example 1: Basic Video Metadata ===")
	url := "https://www.youtube.com/watch?v=jNQXAC9IVRw"
	
	metadata, err := downloader.GetVideoMetadata(url)
	if err != nil {
		log.Fatalf("Failed to get metadata: %v", err)
	}
	
	fmt.Printf("Title: %s\n", metadata.Title)
	fmt.Printf("Uploader: %s\n", metadata.Uploader)
	fmt.Printf("Duration: %s\n", metadata.DurationString)
	fmt.Printf("Views: %d\n", metadata.ViewCount)
	fmt.Printf("Likes: %d\n", metadata.LikeCount)
	fmt.Printf("Upload Date: %s\n", metadata.UploadDate)
	
	// Example 2: Check video quality before downloading
	fmt.Println("\n=== Example 2: Video Quality Information ===")
	fmt.Printf("Resolution: %dx%d\n", metadata.Width, metadata.Height)
	fmt.Printf("FPS: %.0f\n", metadata.FPS)
	fmt.Printf("Video Codec: %s\n", metadata.VideoCodec)
	fmt.Printf("Audio Codec: %s\n", metadata.AudioCodec)
	fmt.Printf("Format: %s\n", metadata.Format)
	
	if metadata.Filesize > 0 {
		fmt.Printf("File Size: %.2f MB\n", float64(metadata.Filesize)/(1024*1024))
	} else if metadata.FilesizeApprox > 0 {
		fmt.Printf("Approx File Size: %.2f MB\n", float64(metadata.FilesizeApprox)/(1024*1024))
	}
	
	// Example 3: Check if video is live
	fmt.Println("\n=== Example 3: Live Status ===")
	if metadata.IsLive {
		fmt.Println("This is a LIVE stream")
	} else if metadata.WasLive {
		fmt.Println("This was a live stream (now archived)")
	} else {
		fmt.Println("This is a regular video")
	}
	fmt.Printf("Live Status: %s\n", metadata.LiveStatus)
	
	// Example 4: Channel Information
	fmt.Println("\n=== Example 4: Channel Information ===")
	fmt.Printf("Channel: %s\n", metadata.Channel)
	fmt.Printf("Channel ID: %s\n", metadata.ChannelID)
	fmt.Printf("Channel URL: %s\n", metadata.ChannelURL)
	
	// Example 5: Additional metadata
	fmt.Println("\n=== Example 5: Additional Info ===")
	if len(metadata.Tags) > 0 {
		fmt.Printf("Tags: %v\n", metadata.Tags)
	}
	if len(metadata.Categories) > 0 {
		fmt.Printf("Categories: %v\n", metadata.Categories)
	}
	fmt.Printf("Thumbnail: %s\n", metadata.Thumbnail)
	
	// Example 6: Decide whether to download based on metadata
	fmt.Println("\n=== Example 6: Smart Download Decision ===")
	
	// Check duration
	maxDurationSeconds := 600 // 10 minutes
	if metadata.Duration > maxDurationSeconds {
		fmt.Printf("Video is %d seconds (>%d seconds), might be too long\n", 
			metadata.Duration, maxDurationSeconds)
	} else {
		fmt.Println("Video duration is acceptable")
	}
	
	// Check file size
	maxSizeMB := 100
	estimatedSize := metadata.FilesizeApprox
	if estimatedSize == 0 {
		estimatedSize = metadata.Filesize
	}
	
	if estimatedSize > 0 {
		sizeMB := estimatedSize / (1024 * 1024)
		if sizeMB > int64(maxSizeMB) {
			fmt.Printf("Video is ~%d MB (>%d MB), might be too large\n", sizeMB, maxSizeMB)
		} else {
			fmt.Printf("Video size (~%d MB) is acceptable, proceeding with download...\n", sizeMB)
			
			// Now download with progress
			progress := func(p downloader.DownloadProgress) {
				fmt.Printf("\rDownloading: %s", p.Stage)
			}
			
			filename, err := downloader.DownloadVideoWithProgress(
				url, "mp4", "720", "avc1", progress,
			)
			if err != nil {
				log.Printf("Download failed: %v", err)
			} else {
				fmt.Printf("\n✓ Downloaded: %s\n", filename)
			}
		}
	}
	
	// Example 7: Access raw metadata for custom fields
	fmt.Println("\n=== Example 7: Raw Metadata Access ===")
	if metadata.Raw != nil {
		// Access any additional fields not in the struct
		if playlist, ok := metadata.Raw["playlist"]; ok {
			fmt.Printf("Playlist: %v\n", playlist)
		}
		if ageLimit, ok := metadata.Raw["age_limit"]; ok {
			fmt.Printf("Age Limit: %v\n", ageLimit)
		}
		fmt.Printf("Total raw fields: %d\n", len(metadata.Raw))
	}
}
