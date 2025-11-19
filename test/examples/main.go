package main

import (
	"fmt"
	"log"
	"time"

	"github.com/debargha2001/gostreampuller/downloader"
)

func main() {
	fmt.Println("🎬 gostreampuller Example Suite")
	fmt.Println("================================")
	fmt.Println()

	// Example 1: Basic video download
	example1BasicVideo()

	// Example 2: Video with progress tracking
	example2VideoWithProgress()

	// Example 3: Audio download
	example3AudioDownload()

	// Example 4: Large file with optimization
	example4LargeFileOptimization()

	fmt.Println("\n✅ All examples completed!")
}

func example1BasicVideo() {
	fmt.Println("📹 Example 1: Basic Video Download")
	fmt.Println("-----------------------------------")
	
	video, err := downloader.DownloadVideo(
		"https://youtube.com/watch?v=jNQXAC9IVRw",
		"mp4",
		"360",
		"avc1",
	)
	if err != nil {
		log.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Downloaded: %s\n\n", video)
}

func example2VideoWithProgress() {
	fmt.Println("📊 Example 2: Video Download with Progress")
	fmt.Println("-------------------------------------------")
	
	progressCallback := func(progress downloader.DownloadProgress) {
		fmt.Printf("[%s] ⏳ %s", time.Now().Format("15:04:05"), progress.Stage)
		if progress.Percentage > 0 {
			fmt.Printf(" - %.1f%%", progress.Percentage)
		}
		fmt.Println()
	}

	video, err := downloader.DownloadVideoWithProgress(
		"https://youtube.com/watch?v=jNQXAC9IVRw",
		"mp4",
		"360",
		"avc1",
		progressCallback,
	)
	if err != nil {
		log.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Downloaded: %s\n\n", video)
}

func example3AudioDownload() {
	fmt.Println("🎵 Example 3: Audio Download")
	fmt.Println("-----------------------------")
	
	progressCallback := func(progress downloader.DownloadProgress) {
		fmt.Printf("⏳ %s...\n", progress.Stage)
	}

	audio, err := downloader.DownloadAudioWithProgress(
		"https://youtube.com/watch?v=jNQXAC9IVRw",
		"mp3",
		"libmp3lame",
		"192k",
		progressCallback,
	)
	if err != nil {
		log.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Downloaded: %s\n\n", audio)
}

func example4LargeFileOptimization() {
	fmt.Println("🚀 Example 4: Large File with Optimization")
	fmt.Println("-------------------------------------------")
	
	// Optimize for large files
	downloader.SetChunkSize(64 * 1024 * 1024)
	downloader.SetMaxConcurrentDownloads(5)
	
	fmt.Println("Settings:")
	fmt.Printf("  Chunk Size: %d MB\n", downloader.ChunkSize/(1024*1024))
	fmt.Printf("  Max Concurrent: %d\n", downloader.MaxConcurrentDownloads)
	fmt.Println()
	
	progressCallback := func(progress downloader.DownloadProgress) {
		if progress.Percentage > 0 {
			fmt.Printf("⏳ %s - %.0f%%\n", progress.Stage, progress.Percentage)
		} else {
			fmt.Printf("⏳ %s...\n", progress.Stage)
		}
	}

	video, err := downloader.DownloadVideoWithProgress(
		"https://youtube.com/watch?v=jNQXAC9IVRw",
		"mp4",
		"720",
		"avc1",
		progressCallback,
	)
	if err != nil {
		log.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Downloaded: %s\n", video)
}
