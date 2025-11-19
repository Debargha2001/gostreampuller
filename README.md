# 🎬 gostreampuller

![gostreampuller banner](https://github.com/user-attachments/assets/6fd472c1-953f-4b0f-97d3-00c0a0e99424)

**gostreampuller** is a lightweight Go package to download videos or extract audio from popular streaming platforms including **YouTube, Facebook, Instagram, TikTok, Twitter, and many more** using [`yt-dlp`](https://github.com/yt-dlp/yt-dlp) and [`ffmpeg`](https://ffmpeg.org/).  
It's designed for cross-platform use: **Linux, macOS, and Windows**.

## ✨ What's New in v1.1.0

- 🎯 **Smart Auto-Installation**: Binaries auto-install on first use (unless you prefer CLI setup)
- 🚀 **Large File Support**: Efficiently handle videos of any size (30+ minutes, 4K, etc.)
- 🔄 **Streaming Architecture**: Uses streaming and chunked processing to prevent memory overflow
- ⚡ **Concurrent Downloads**: Download video fragments concurrently for better performance
- 📊 **Progress Tracking**: Optional callbacks to monitor download and conversion progress
- 📋 **Video Metadata**: Fetch comprehensive video information without downloading
- ⏱️ **Timeout Protection**: Automatic timeout handling with configurable durations
- 🔧 **Customizable Settings**: Adjust chunk size and concurrency for optimal performance
- 🛡️ **Retry Logic**: Built-in retry mechanism for failed downloads and fragments
- 🎯 **CLI Tool Option**: Optional CLI tool for manual dependency management
- 🧪 **Production-Ready Testing**: Comprehensive test suite in dedicated `test/` directory

---

## 📦 Installation

### Zero Configuration (Recommended - New in v1.1.0! 🎉)

Just install and use - binaries are **automatically installed on first use**:

```bash
# Install the package
go get github.com/debargha2001/gostreampuller
```

That's it! On your first download/metadata call, the package will:
1. Detect if binaries are needed
2. Auto-install yt-dlp and ffmpeg to `~/.gostreampuller/bin/` (takes 1-3 min, one time)
3. Show clear progress messages
4. Subsequent calls are instant!

```go
import "github.com/debargha2001/gostreampuller/downloader"

func main() {
    // First call auto-installs binaries if needed (one-time setup)
    // After that, it's instant!
    video, _ := downloader.DownloadVideo(url, "mp4", "720", "avc1")
}
```

**Note**: Auto-installation is smart:
- ✅ Installs automatically if you never used CLI setup
- ✅ Respects system binaries (`brew install`, `apt install`, etc.)
- ✅ Skipped if you prefer manual CLI setup (see below)

---

### Manual CLI Setup (Optional)

If you prefer explicit control over installation:

```bash
# Install the CLI tool
go install github.com/debargha2001/gostreampuller/cmd/gostreampuller-cli@latest

# Run setup manually
gostreampuller-cli setup

# Check installation status
gostreampuller-cli check
```

When you use CLI setup, auto-installation is automatically disabled to respect your choice.

---

### Disable Auto-Installation

For corporate environments or custom setups:

```bash
# Completely disable auto-installation
export GOSTREAMPULLER_NO_AUTO_INSTALL=1

# Then install manually
gostreampuller-cli setup
# OR use system package managers
brew install yt-dlp ffmpeg  # macOS
sudo apt install yt-dlp ffmpeg  # Ubuntu/Debian
```

---

## 🚀 Usage

### Zero Configuration Usage (New in v1.1.0!)

Just import and use - binaries install automatically on first call:

```go
package main

import (
  "fmt"
  "log"

  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  // First call? Binaries auto-install (1-3 min, one time, clear progress shown)
  // Already installed? Works instantly!
  
  // Returns the absolute path to the downloaded file
  video, err := downloader.DownloadVideo("https://youtube.com/watch?v=dQw4w9WgXcQ", "mp4", "720", "avc1")
  if err != nil {
    log.Fatal(err)
  }
  // video = "/full/path/to/current/directory/video_1234567890.mp4"
  fmt.Println("Downloaded video:", video)

  // Subsequent calls are instant
  audio, err := downloader.DownloadAudio("https://youtube.com/watch?v=example", "mp3", "libmp3lame", "128k")
  if err != nil {
    log.Fatal(err)
  }
  // audio = "/full/path/to/current/directory/audio_1234567891.mp3"
  fmt.Println("Downloaded audio:", audio)
}
```

**📁 Where are files saved?**
- **Location**: Current working directory by default (where your program runs)
- **Custom location**: Use `DownloadVideoToDir()` or `DownloadAudioToDir()` to specify a custom directory
- **Naming**: `video_<timestamp>.<format>` or `audio_<timestamp>.<format>`
- **Return value**: Absolute path to the downloaded file

**Example:**
```bash
$ cd /home/user/downloads
$ go run main.go
# Files saved in: /home/user/downloads/
# - video_1732021234567890.mp4
# - audio_1732021234567891.mp3
```

**What you'll see on first run:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  gostreampuller: First-time setup
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Required binaries are being installed automatically.
  This is a one-time process and takes 1-3 minutes.
  
  To disable: export GOSTREAMPULLER_NO_AUTO_INSTALL=1
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[gostreampuller] Installing yt-dlp...
[gostreampuller] ✓ yt-dlp installed successfully
[gostreampuller] Installing ffmpeg...
[gostreampuller] ✓ ffmpeg installed successfully

[gostreampuller] ✓ Setup complete! Subsequent calls will be fast.

Downloaded video: video.mp4
```

---

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
  "fmt"
  "log"

  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  url := "https://youtube.com/watch?v=dQw4w9WgXcQ"
  
  // Download video to current directory
  video, err := downloader.DownloadVideo(url, "mp4", "720", "avc1")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Video:", video)
  
  // Download to custom directory
  video, err = downloader.DownloadVideoToDir(url, "mp4", "1080", "avc1", "/downloads")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Video saved to:", video)
}
```

### API Overview

**Video Downloads:**
- `DownloadVideo()` - Download to current directory
- `DownloadVideoToDir()` - Download to custom directory
- `DownloadVideoWithProgress()` - With progress tracking
- `DownloadVideoToDirWithProgress()` - Custom directory + progress

**Audio Downloads:**
- `DownloadAudio()` - Download to current directory  
- `DownloadAudioToDir()` - Download to custom directory
- `DownloadAudioWithProgress()` - With progress tracking
- `DownloadAudioToDirWithProgress()` - Custom directory + progress

**Metadata:**
- `GetVideoMetadata()` - Fetch video information without downloading
- `GetVideoMetadataWithContext()` - With timeout control

**Binary Management:**
- `SetYTDLPPath()` - Override yt-dlp binary location
- `SetFFMPEGPath()` - Override ffmpeg binary location
- `ResetBinaryPaths()` - Reset to auto-detected paths

---

### Basic Usage (Backward Compatible)

```go
package main

import (
  "fmt"
  "log"

  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  // Auto-installation handles binaries automatically
  // But you can still override paths if needed (optional):
  // downloader.SetYTDLPPath("C:\\Users\\yourname\\yt-dlp.exe")
  // downloader.SetFFMPEGPath("C:\\ffmpeg\\bin\\ffmpeg.exe")

  // Basic video download with defaults (mp4, 720p, avc1 codec)
  video, err := downloader.DownloadVideo("https://youtube.com/watch?v=dQw4w9WgXcQ", "", "", "")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Downloaded video:", video)  // Prints absolute path

  // Basic audio download with defaults (mp3, libmp3lame codec, 128k bitrate)
  audio, err := downloader.DownloadAudio("https://instagram.com/p/example", "", "", "")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Downloaded audio:", audio)  // Prints absolute path
}
```

---

### Specifying Custom Output Directory

Use the `*ToDir` functions to download directly to a specific directory:

```go
package main

import (
  "fmt"
  "log"

  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  // Download video to specific directory
  video, err := downloader.DownloadVideoToDir(
    "https://youtube.com/watch?v=example",
    "mp4",
    "1080",
    "avc1",
    "/path/to/downloads", // Directory will be created if it doesn't exist
  )
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Video saved to: %s\n", video)
  // Output: Video saved to: /path/to/downloads/video_1732021234567890.mp4
  
  // Download audio to different directory
  audio, err := downloader.DownloadAudioToDir(
    "https://youtube.com/watch?v=example",
    "mp3",
    "libmp3lame",
    "320k",
    "/music/downloads",
  )
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Audio saved to: %s\n", audio)
  // Output: Audio saved to: /music/downloads/audio_1732021234567891.mp3
  
  // Use empty string "" for current directory (same as DownloadVideo/DownloadAudio)
  video2, err := downloader.DownloadVideoToDir(url, "mp4", "720", "avc1", "")
  // Saves to current directory
}
```

**With Progress Tracking:**

```go
progressCallback := func(progress downloader.DownloadProgress) {
  fmt.Printf("[%s] %.1f%%\n", progress.Stage, progress.Percentage)
}

video, err := downloader.DownloadVideoToDirWithProgress(
  url, "mp4", "1080", "avc1",
  "/downloads", // Custom directory
  progressCallback,
)
```

---

### Large File Downloads with Progress Tracking

```go
package main

import (
  "fmt"
  "log"
  "time"

  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  // Progress callback to monitor download status
  progressCallback := func(progress downloader.DownloadProgress) {
    fmt.Printf("[%s] %s", time.Now().Format("15:04:05"), progress.Stage)
    if progress.Percentage > 0 {
      fmt.Printf(" - %.1f%%", progress.Percentage)
    }
    fmt.Println()
  }

  // Download large video (30+ minutes, 4K, etc.) with progress tracking
  video, err := downloader.DownloadVideoWithProgress(
    "https://youtube.com/watch?v=long_video_example",
    "mp4",
    "2160", // 4K resolution
    "avc1",
    progressCallback,
  )
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Downloaded large video:", video)

  // Download large audio file (podcasts, long music, etc.) with progress
  audio, err := downloader.DownloadAudioWithProgress(
    "https://youtube.com/watch?v=podcast_1hour",
    "mp3",
    "libmp3lame",
    "320k",
    progressCallback,
  )
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Downloaded large audio:", audio)
}
```

### Optimizing for Very Large Files (New in v1.1.0)

```go
package main

import (
  "fmt"
  "log"

  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  // Customize chunk size for better performance (default: 32MB)
  downloader.SetChunkSize(64 * 1024 * 1024) // 64MB chunks for very large files

  // Increase concurrent downloads (default: 3)
  downloader.SetMaxConcurrentDownloads(5)

  // Now downloads will use these optimized settings
  video, err := downloader.DownloadVideo(
    "https://youtube.com/watch?v=4k_2hour_movie",
    "mp4",
    "2160",
    "avc1",
  )
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Downloaded optimized video:", video)
}
```

### Advanced Usage with Custom Parameters

```go
package main

import (
  "fmt"
  "log"

  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  // Download 1080p video in webm format with VP9 codec
  video, err := downloader.DownloadVideo(
    "https://facebook.com/watch?v=example",
    "webm",  // format
    "1080",  // resolution
    "vp9",   // codec
  )
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Downloaded HD video:", video)

  // Download high-quality audio in FLAC format
  audio, err := downloader.DownloadAudio(
    "https://tiktok.com/@user/video/example",
    "flac",      // output format
    "flac",      // codec
    "320k",      // bitrate
  )
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Downloaded FLAC audio:", audio)
}
```

### Using Custom Binary Paths (Optional)

By default, the package intelligently manages binaries:
1. **Auto-installs** to `~/.gostreampuller/bin/` on first use (if not using CLI)
2. **Uses CLI-installed** binaries if you ran `gostreampuller-cli setup`
3. **Detects system binaries** from PATH (`brew install`, `apt install`, etc.)

You can override these defaults if needed:

```go
package main

import (
  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  // Override with custom paths (optional)
  downloader.SetYTDLPPath("/usr/local/bin/yt-dlp")
  downloader.SetFFMPEGPath("/usr/local/bin/ffmpeg")
  
  // Or on Windows:
  // downloader.SetYTDLPPath("C:\\tools\\yt-dlp.exe")
  // downloader.SetFFMPEGPath("C:\\ffmpeg\\bin\\ffmpeg.exe")
  
  // Now downloads will use your custom paths
  video, _ := downloader.DownloadVideo(url, "mp4", "720", "avc1")
  
  // To revert back to auto-detected binaries:
  downloader.ResetBinaryPaths()
}
```

**Note**: You rarely need to set custom paths. The smart auto-detection handles most scenarios automatically.

---

## 🎯 Smart Binary Management

### How It Works

The package intelligently manages dependencies with **zero configuration required**:

```
┌─────────────────────────────────────────────────────┐
│  Scenario 1: Fresh Install (Auto-Install)          │
├─────────────────────────────────────────────────────┤
│  First API call:                                    │
│    → Check: Binaries exist? NO                     │
│    → Check: System binaries? NO                    │
│    → Action: Auto-install to ~/.gostreampuller/bin │
│    → Time: 1-3 min (one time only)                 │
│    → Result: Subsequent calls instant ✓            │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│  Scenario 2: Manual CLI Setup                      │
├─────────────────────────────────────────────────────┤
│  User runs: gostreampuller-cli setup                │
│    → Installs binaries + creates marker file       │
│  API calls:                                         │
│    → Uses CLI-installed binaries                   │
│    → NO auto-install (respects user choice) ✓      │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│  Scenario 3: System Package Manager                │
├─────────────────────────────────────────────────────┤
│  User installs: brew install yt-dlp ffmpeg          │
│  API calls:                                         │
│    → Detects system binaries in PATH               │
│    → Uses system binaries                          │
│    → NO auto-install (respects system install) ✓   │
└─────────────────────────────────────────────────────┘
```

### Environment Variables

```bash
# Disable auto-installation completely
export GOSTREAMPULLER_NO_AUTO_INSTALL=1

# Show verbose installation progress
export GOSTREAMPULLER_VERBOSE=1
```

### Priority Order

The package checks for binaries in this order:

1. **Custom paths** (if set via `SetYTDLPPath()`/`SetFFMPEGPath()`)
2. **Local binaries** (`~/.gostreampuller/bin/`)
3. **System binaries** (from PATH)
4. **Auto-install** (if none found and not disabled)

---

## 📋 Fetching Video Metadata (New in v1.1.0)

Get comprehensive video information without downloading:

```go
package main

import (
  "fmt"
  "log"

  "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
  // Fetch video metadata
  metadata, err := downloader.GetVideoMetadata("https://youtube.com/watch?v=dQw4w9WgXcQ")
  if err != nil {
    log.Fatal(err)
  }

  // Access video information
  fmt.Printf("Title: %s\n", metadata.Title)
  fmt.Printf("Uploader: %s\n", metadata.Uploader)
  fmt.Printf("Duration: %s (%d seconds)\n", metadata.DurationString, metadata.Duration)
  fmt.Printf("Views: %d\n", metadata.ViewCount)
  fmt.Printf("Likes: %d\n", metadata.LikeCount)
  fmt.Printf("Resolution: %dx%d\n", metadata.Width, metadata.Height)
  fmt.Printf("Video Codec: %s\n", metadata.VideoCodec)
  fmt.Printf("Audio Codec: %s\n", metadata.AudioCodec)
  fmt.Printf("Upload Date: %s\n", metadata.UploadDate)
  
  // Check if live
  if metadata.IsLive {
    fmt.Println("This is a LIVE stream")
  }
  
  // Check file size before downloading
  if metadata.FilesizeApprox > 0 {
    sizeMB := float64(metadata.FilesizeApprox) / (1024 * 1024)
    fmt.Printf("Estimated Size: %.2f MB\n", sizeMB)
    
    // Decide whether to download based on size
    if sizeMB < 100 {
      video, _ := downloader.DownloadVideo(metadata.WebpageURL, "mp4", "720", "avc1")
      fmt.Println("Downloaded:", video)
    } else {
      fmt.Println("Video too large, skipping download")
    }
  }
  
  // Access raw metadata for additional fields
  if ageLimit, ok := metadata.Raw["age_limit"]; ok {
    fmt.Printf("Age Limit: %v\n", ageLimit)
  }
}
```

### VideoMetadata Structure

The `VideoMetadata` struct contains comprehensive information:

- **Basic**: `ID`, `Title`, `Description`, `Uploader`, `UploaderID`, `UploaderURL`
- **Stats**: `Duration`, `ViewCount`, `LikeCount`, `DislikeCount`, `AverageRating`
- **Technical**: `Width`, `Height`, `FPS`, `VideoCodec`, `AudioCodec`, `Format`, `Filesize`
- **URLs**: `URL`, `WebpageURL`, `Thumbnail`
- **Dates**: `UploadDate`, `ReleaseDate`, `Timestamp`
- **Channel**: `Channel`, `ChannelID`, `ChannelURL`
- **Live**: `IsLive`, `WasLive`, `LiveStatus`
- **Additional**: `Categories`, `Tags`, `Subtitles`
- **Raw**: Access to all original fields from yt-dlp

---

## 📋 Function Parameters

### GetVideoMetadata Parameters

```go
func GetVideoMetadata(url string) (*VideoMetadata, error)
func GetVideoMetadataWithContext(ctx context.Context, url string) (*VideoMetadata, error)
```

| Parameter | Description | Required |
|-----------|-------------|----------|
| `url` | Video URL from supported platforms | Yes |

**Returns**: `*VideoMetadata` with comprehensive video information or error

### DownloadVideo Parameters

```go
// Download to current directory
func DownloadVideo(url string, format string, resolution string, codec string) (string, error)
func DownloadVideoWithProgress(url string, format string, resolution string, codec string, progressCb ProgressCallback) (string, error)

// Download to custom directory
func DownloadVideoToDir(url string, format string, resolution string, codec string, outputDir string) (string, error)
func DownloadVideoToDirWithProgress(url string, format string, resolution string, codec string, outputDir string, progressCb ProgressCallback) (string, error)
```

| Parameter | Description | Default | Possible Values |
|-----------|-------------|---------|-----------------|
| `url` | Video URL from supported platforms | *required* | YouTube, Facebook, Instagram, TikTok, Twitter, etc. |
| `format` | Output video format | `mp4` | `mp4`, `webm`, `mkv`, `avi`, `mov`, `flv` |
| `resolution` | Video resolution height | `720` | `144`, `240`, `360`, `480`, `720`, `1080`, `1440`, `2160` (4K) |
| `codec` | Video codec | `avc1` | `avc1` (H.264), `vp9`, `vp8`, `av01` (AV1), `hevc` (H.265) |
| `outputDir` | Custom output directory (*ToDir only) | *current dir* | Any path. Created if doesn't exist. Use `""` for current directory. |

**Returns**: `string` - Absolute path to the downloaded video file

**File Location**: 
- `DownloadVideo*`: Current working directory → `video_<timestamp>.<format>`
- `DownloadVideoToDir*`: Specified directory → `<outputDir>/video_<timestamp>.<format>`

**Example combinations:**
```go
// HD MP4 to current directory
video, _ := downloader.DownloadVideo(url, "mp4", "1080", "avc1")
// Returns: "/current/dir/video_1732021234567890.mp4"

// 4K WebM to custom directory
video, _ := downloader.DownloadVideoToDir(url, "webm", "2160", "vp9", "/downloads")
// Returns: "/downloads/video_1732021234567891.webm"

// With progress tracking to custom directory
progressCb := func(p downloader.DownloadProgress) {
  fmt.Printf("[%s] %.1f%%\n", p.Stage, p.Percentage)
}
video, _ := downloader.DownloadVideoToDirWithProgress(url, "mp4", "1080", "avc1", "/videos", progressCb)
// Returns: "/videos/video_1732021234567892.mp4"
```

### DownloadAudio Parameters

```go
// Download to current directory
func DownloadAudio(url string, outputFormat string, codec string, bitrate string) (string, error)
func DownloadAudioWithProgress(url string, outputFormat string, codec string, bitrate string, progressCb ProgressCallback) (string, error)

// Download to custom directory
func DownloadAudioToDir(url string, outputFormat string, codec string, bitrate string, outputDir string) (string, error)
func DownloadAudioToDirWithProgress(url string, outputFormat string, codec string, bitrate string, outputDir string, progressCb ProgressCallback) (string, error)
```

| Parameter | Description | Default | Possible Values |
|-----------|-------------|---------|-----------------|
| `url` | Video URL from supported platforms | *required* | YouTube, Facebook, Instagram, TikTok, Twitter, etc. |
| `outputFormat` | Output audio format | `mp3` | `mp3`, `aac`, `ogg`, `wav`, `flac`, `m4a` |
| `codec` | Audio codec | `libmp3lame` | `libmp3lame`, `aac`, `libvorbis`, `pcm_s16le`, `flac`, `libfdk_aac` |
| `bitrate` | Audio bitrate | `128k` | `64k`, `96k`, `128k`, `192k`, `256k`, `320k` |
| `outputDir` | Custom output directory (*ToDir only) | *current dir* | Any path. Created if doesn't exist. Use `""` for current directory. |

**Returns**: `string` - Absolute path to the downloaded audio file

**File Location**: 
- `DownloadAudio*`: Current working directory → `audio_<timestamp>.<format>`
- `DownloadAudioToDir*`: Specified directory → `<outputDir>/audio_<timestamp>.<format>`

**Example combinations:**
```go
// High-quality MP3 to current directory
audio, _ := downloader.DownloadAudio(url, "mp3", "libmp3lame", "320k")
// Returns: "/current/dir/audio_1732021234567890.mp3"

// Lossless FLAC to custom directory
audio, _ := downloader.DownloadAudioToDir(url, "flac", "flac", "320k", "/music")
// Returns: "/music/audio_1732021234567891.flac"

// Efficient AAC with progress tracking
progressCb := func(p downloader.DownloadProgress) {
  fmt.Printf("[%s] %.1f%%\n", p.Stage, p.Percentage)
}
audio, _ := downloader.DownloadAudioToDirWithProgress(url, "m4a", "aac", "192k", "/audio", progressCb)
// Returns: "/audio/audio_1732021234567892.m4a"
```

### Parameter Combinations Guide

**For highest quality video:**
```go
downloader.DownloadVideo(url, "webm", "2160", "vp9")  // 4K WebM with VP9
```

**For best compatibility:**
```go
downloader.DownloadVideo(url, "mp4", "1080", "avc1")  // 1080p MP4 with H.264
```

**For smallest file size:**
```go
downloader.DownloadVideo(url, "webm", "720", "vp9")   // 720p WebM with VP9
```

**For highest quality audio:**
```go
downloader.DownloadAudio(url, "flac", "flac", "320k") // Lossless FLAC
```

**For best compatibility audio:**
```go
downloader.DownloadAudio(url, "mp3", "libmp3lame", "320k") // High-quality MP3
```

**For smallest audio file:**
```go
downloader.DownloadAudio(url, "ogg", "libvorbis", "128k") // Efficient OGG Vorbis
```

---

## 🛠️ Platform Setup

### Automated Setup (Easiest - All Platforms)

The package now includes a CLI tool that automatically installs yt-dlp and ffmpeg for you:

```bash
# Install CLI tool
go install github.com/debargha2001/gostreampuller/cmd/gostreampuller-cli@latest

# Run automated setup
gostreampuller-cli setup

# Verify installation
gostreampuller-cli check
```

**Benefits:**
- ✅ Works on Linux, macOS, and Windows
- ✅ No manual PATH configuration needed
- ✅ Binaries stored in `~/.gostreampuller/bin/`
- ✅ Package auto-detects and uses local binaries
- ✅ Perfect for CI/CD and offline environments

### Manual Installation (Alternative)

If you prefer to install manually or already have yt-dlp/ffmpeg installed system-wide, they'll work automatically.

### 🐧 Linux Setup

**Option 1: Using package managers**

**Ubuntu/Debian:**
```bash
# Install yt-dlp
sudo apt update
sudo apt install yt-dlp

# Install ffmpeg
sudo apt install ffmpeg
```

**Arch Linux:**
```bash
# Install yt-dlp
sudo pacman -S yt-dlp

# Install ffmpeg
sudo pacman -S ffmpeg
```

**Option 2: Manual installation**
```bash
# Download yt-dlp
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp

# Install ffmpeg (if not available via package manager)
# Download from https://ffmpeg.org/download.html
```

### 🍎 macOS Setup

**Option 1: Using Homebrew (recommended)**
```bash
# Install Homebrew if you haven't already
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install yt-dlp
brew install yt-dlp

# Install ffmpeg
brew install ffmpeg
```

**Option 2: Using MacPorts**
```bash
# Install yt-dlp
sudo port install yt-dlp

# Install ffmpeg
sudo port install ffmpeg
```

**Option 3: Manual installation**
```bash
# Download yt-dlp
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp

# For ffmpeg, download from https://ffmpeg.org/download.html
```

### 🪟 Windows Setup

**Option 1: Using package managers**

**Using Chocolatey:**
```powershell
# Install Chocolatey if you haven't already
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install yt-dlp
choco install yt-dlp

# Install ffmpeg
choco install ffmpeg
```

**Using Scoop:**
```powershell
# Install Scoop if you haven't already
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# Install yt-dlp
scoop install yt-dlp

# Install ffmpeg
scoop install ffmpeg
```

**Option 2: Manual installation**

1. **Download binaries**:
   - [yt-dlp.exe](https://github.com/yt-dlp/yt-dlp/releases/latest)
   - [ffmpeg Windows build](https://ffmpeg.org/download.html)

2. **Add to PATH or override in Go**:

```go
// If not in PATH, override binary paths in Go
downloader.SetYTDLPPath("C:\\Users\\yourname\\Desktop\\yt-dlp.exe")
downloader.SetFFMPEGPath("C:\\ffmpeg\\bin\\ffmpeg.exe")
```

### 🔧 Custom Binary Paths

If you prefer to use custom binary locations on any platform, you can override the default paths:

```go
// Linux/macOS examples
downloader.SetYTDLPPath("/usr/local/bin/yt-dlp")
downloader.SetFFMPEGPath("/usr/local/bin/ffmpeg")

// Windows example
downloader.SetYTDLPPath("C:\\Users\\yourname\\Desktop\\yt-dlp.exe")
downloader.SetFFMPEGPath("C:\\ffmpeg\\bin\\ffmpeg.exe")
```

That's it! ✅

---

## 🐳 Docker Support

For containerized environments, you can use this Dockerfile to create a container with gostreampuller, yt-dlp, and ffmpeg pre-installed:

```dockerfile
# ---------- Build Stage ----------
# Use official Golang image to compile the Go app
FROM golang:1.24.1-bullseye AS builder

# Set the working directory inside the builder container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code into the builder container
COPY . .

# Build the Go binary; replace 'your-app-name' with your actual binary name
RUN go build -o your-app-name main.go

# ---------- Runtime Stage ----------
# Use a lightweight Debian image for production
FROM debian:bullseye-slim

# Install required tools:
# - python3: needed by yt-dlp
# - ffmpeg: for audio processing
# - curl: to download yt-dlp binary
# - ca-certificates: ensures HTTPS requests work
RUN apt-get update && \
    apt-get install -y python3 ffmpeg curl ca-certificates && \
    curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the runtime container
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/your-app-name .

# Expose the API port (adjust port as needed)
EXPOSE 8080

# Run the Go server when the container starts
CMD ["./your-app-name"]
```

> **📝 Note**: Replace `your-app-name` with the actual name of your Go application binary. The Dockerfile above is a template that can be customized for any Go application using gostreampuller.

### 🔨 Building and Running with Docker

```bash
# Build the Docker image
docker build -t gostreampuller-app .

# Run the container
docker run -p 8080:8080 gostreampuller-app
```

**Benefits of using Docker:**
- ✅ No need to manually install yt-dlp and ffmpeg
- ✅ Consistent environment across different systems
- ✅ Easy deployment and scaling
- ✅ Isolated dependencies

---

## 🧠 Features

### Core Functionality
- 📹 Download videos in multiple formats (`.mp4`, `.webm`, `.mkv`, etc.)
- 🎧 Extract audio in various formats (`.mp3`, `.flac`, `.aac`, etc.)
- 🌐 Supports multiple platforms: **YouTube, Facebook, Instagram, TikTok, Twitter, and more**
- 🧩 Works across Windows, macOS, Linux
- 🧼 Automatically cleans up temporary files
- ⚡ Smart format detection and conversion

### v1.1.0 Enhancements
- 🎯 **Smart Auto-Installation** - Binaries install automatically on first use (respects CLI/system installs)
- 🚀 **Large File Support** - Handles 30+ min videos, 4K content efficiently
- 🔄 **Streaming Architecture** - Prevents memory overflow with chunked processing
- 📊 **Progress Tracking** - Real-time callbacks for download status
- ⚡ **Concurrent Downloads** - Parallel fragment downloads for better performance
- 🛡️ **Retry Logic** - Automatic retries with configurable timeouts
- 🔧 **Customizable Settings** - Adjust chunk size and concurrency
- 📋 **Video Metadata** - Fetch comprehensive video info without downloading
- 🎯 **Optional CLI Tool** - Manual dependency management for those who prefer it
- 🧪 **Production-Ready Tests** - Comprehensive test suite with automated runner

---

## 🗂️ Folder Structure

```
gostreampuller/
├── cmd/
│   └── gostreampuller-cli/     // CLI tool for setup
│       └── main.go
├── downloader/
│   └── downloader.go            // Core package with streaming support
├── internal/
│   └── installer/               // Automated installer for dependencies
│       └── installer.go
├── test/                        // Comprehensive test suite
│   ├── integration_test.go      // Integration tests
│   ├── unit_test.go             // Unit tests
│   ├── run_tests.sh             // Test runner script
│   ├── README.md                // Test documentation
│   └── examples/                // Usage examples
│       └── main.go
├── go.mod
└── README.md
```

---

## 🧪 Testing

The package includes a comprehensive test suite in the `test/` directory with production-grade testing infrastructure.

### Quick Start

```bash
# Run all tests with the automated test runner
cd test
./run_tests.sh

# Run with all options (integration, examples, benchmarks, coverage)
./run_tests.sh --integration --examples --bench --coverage
# or use short flags
./run_tests.sh -i -e -b -c
```

### Test Structure

```
test/
├── run_tests.sh           # Main test runner (executable)
├── full_system_test.sh    # Interactive comprehensive tests
├── integration_test.go    # Real download tests
├── unit_test.go          # Fast unit tests
├── metadata_test.go      # Metadata fetching tests
├── README.md             # Testing documentation
└── examples/
    ├── main.go           # Basic usage examples
    ├── metadata/         # Metadata examples
    └── custom_directory/ # Custom directory examples
```

### Test Categories

**Unit Tests** (< 1 minute):
```bash
go test -v -short ./test
```
- Package initialization
- Configuration setters
- Type validation
- No downloads required

**Integration Tests** (1-30 minutes):
```bash
go test -v ./test -timeout 30m
```
- Real video/audio downloads
- Progress tracking
- Concurrent operations
- Requires internet + binaries

**Metadata Tests** (10-30 seconds):
```bash
go test -v -run TestGetVideoMetadata ./test
```
- Fetch video information
- Fast, no video data downloaded
- Requires internet + yt-dlp

**Full System Test** (Interactive):
```bash
cd test
./full_system_test.sh
```
- Complete functionality testing
- CLI tool verification
- Binary installation check
- All features tested

### Test Options

```bash
./run_tests.sh                    # Unit tests only (fast)
./run_tests.sh --integration      # Include integration tests
./run_tests.sh --examples         # Run example programs
./run_tests.sh --bench            # Run benchmarks
./run_tests.sh --coverage         # Generate coverage report
./run_tests.sh -i -e -b -c        # Run everything
```

### Manual Testing Commands

```bash
# Unit tests (fast, no downloads)
go test -v -short ./test

# Integration tests (requires internet)
go test -v ./test -timeout 30m

# Specific test
go test -v ./test -run TestDownloadVideo

# Metadata only
go test -v -run TestGetVideoMetadata ./test

# With coverage
go test -coverprofile=coverage.out ./test
go tool cover -html=coverage.out
```

### Running Examples

```bash
# Basic usage examples
cd test/examples && go run main.go

# Metadata examples
cd test/examples/metadata && go run main.go

# Custom directory examples
cd test/examples/custom_directory && go run main.go
```

### Testing Workflow

**For Development:**
```bash
# 1. Quick validation
go build ./...
go test -v -short ./test

# 2. If working on metadata
go test -v -run TestGetVideoMetadata ./test

# 3. If working on downloads
go test -v -run TestDownloadVideo ./test -timeout 10m

# 4. Before committing
./test/full_system_test.sh
```

**For CI/CD:**
```bash
# Install dependencies
gostreampuller-cli setup

# Run all tests
go test -v ./test -timeout 30m

# Run with coverage
go test -v ./test -timeout 30m -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**For Production Validation:**
```bash
# Full comprehensive test
cd test
./full_system_test.sh

# Answer "y" to all prompts to test everything
```

### Debugging Failed Tests

**If binary installation fails:**
```bash
# Check manually
gostreampuller-cli check

# Check installation directory
ls -la ~/.gostreampuller/bin/

# Try manual install
gostreampuller-cli setup
```

**If metadata tests fail:**
```bash
# Test yt-dlp directly
~/.gostreampuller/bin/yt-dlp --dump-json "https://www.youtube.com/watch?v=jNQXAC9IVRw"

# Check network
curl -I https://www.youtube.com
```

**If download tests fail:**
```bash
# Test binaries directly
~/.gostreampuller/bin/yt-dlp --version
~/.gostreampuller/bin/ffmpeg -version

# Check disk space
df -h .

# Run single test with timeout
go test -v -run "TestDownloadVideo$" ./test -timeout 15m
```

### Test Checklist

Before release, verify:

- [ ] CLI tool builds: `go build ./cmd/gostreampuller-cli`
- [ ] Binary installation works: `gostreampuller-cli setup`
- [ ] Binaries are detected: `gostreampuller-cli check`
- [ ] Package builds: `go build ./...`
- [ ] Unit tests pass: `go test -short ./test`
- [ ] Metadata fetching works: `go test -run TestGetVideoMetadata ./test`
- [ ] Video download works: Test with short video
- [ ] Audio download works: Test with short video
- [ ] Progress tracking works: Callbacks are called
- [ ] Custom paths work: `SetYTDLPPath()` / `SetFFMPEGPath()`
- [ ] Path reset works: `ResetBinaryPaths()`
- [ ] Examples run: Both basic and metadata examples
- [ ] Full system test passes: `./test/full_system_test.sh`

See [`test/README.md`](test/README.md) for detailed testing documentation.

---

---

## 📚 Migration Guide: v1.0.0 to v1.1.0

### Summary

Version 1.1.0 is **100% backward compatible** with v1.0.0. All existing code will continue to work without any changes. This guide explains new features and how to take advantage of them.

### What's Changed?

#### ✅ Backward Compatible (No Changes Required)

All v1.0.0 functions work exactly as before:

```go
// These still work perfectly - no changes needed!
video, err := downloader.DownloadVideo(url, "mp4", "720", "avc1")
audio, err := downloader.DownloadAudio(url, "mp3", "libmp3lame", "128k")
```

#### ✨ New Features Available

**1. Auto-Installation (Automatic)**

Your existing code now benefits from smart binary management:
- Automatically installs binaries on first use
- Respects CLI setup if used
- Detects system binaries
- **No code changes needed** - it just works!

**2. Large File Support (Automatic)**

Improvements happen automatically:
- Streaming architecture (no memory overflow)
- Concurrent fragment downloads
- Automatic retry logic
- Better timeout handling
- **No code changes needed** - better performance automatically!

**3. Custom Output Directories (New Option)**

```go
// Before (v1.0.0) - saves to current directory
video, _ := downloader.DownloadVideo(url, "mp4", "1080", "avc1")

// After (v1.1.0) - optional custom directory
video, _ := downloader.DownloadVideoToDir(url, "mp4", "1080", "avc1", "/downloads")
```

**4. Progress Tracking (New Option)**

```go
// Before (v1.0.0) - no progress visibility
video, err := downloader.DownloadVideo(url, "mp4", "1080", "avc1")

// After (v1.1.0) - optional progress tracking
progressCallback := func(progress downloader.DownloadProgress) {
    fmt.Printf("Stage: %s (%.1f%%)\n", progress.Stage, progress.Percentage)
}

video, err := downloader.DownloadVideoWithProgress(
    url, "mp4", "1080", "avc1",
    progressCallback,
)
```

**5. Video Metadata (New Feature)**

```go
// New in v1.1.0 - fetch video info without downloading
metadata, err := downloader.GetVideoMetadata(url)
fmt.Printf("Title: %s\n", metadata.Title)
fmt.Printf("Duration: %s\n", metadata.DurationString)
fmt.Printf("Size: %.2f MB\n", float64(metadata.FilesizeApprox)/(1024*1024))

// Smart download decision
if metadata.Duration < 600 {
    video, _ := downloader.DownloadVideo(url, "mp4", "720", "avc1")
}
```

**6. Performance Tuning (New Option)**

```go
// Before (v1.0.0) - no tuning options
video, err := downloader.DownloadVideo(url, "mp4", "2160", "avc1")

// After (v1.1.0) - optional performance tuning for large files
downloader.SetChunkSize(64 * 1024 * 1024)        // 64MB chunks
downloader.SetMaxConcurrentDownloads(5)          // More concurrency

video, err := downloader.DownloadVideo(url, "mp4", "2160", "avc1")
```

### Migration Checklist

#### For Most Users: Do Nothing! ✅

If your current code works, it will continue to work and perform better automatically.

#### If You Want New Features:

1. **Custom directories**: Use `DownloadVideoToDir()` / `DownloadAudioToDir()`
2. **Progress tracking**: Use `*WithProgress()` variants
3. **Metadata**: Use `GetVideoMetadata()` before downloading
4. **Performance tuning**: Call `SetChunkSize()` / `SetMaxConcurrentDownloads()`

### Common Upgrade Scenarios

#### Scenario 1: Basic Usage (No Changes)

```go
// v1.0.0 and v1.1.0 - same code!
import "github.com/debargha2001/gostreampuller/downloader"

func main() {
    video, err := downloader.DownloadVideo(
        "https://youtube.com/watch?v=example",
        "mp4", "720", "avc1",
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Downloaded:", video)
}
```

#### Scenario 2: Adding Custom Directory

```go
// v1.0.0
func downloadVideo(url string) (string, error) {
    return downloader.DownloadVideo(url, "mp4", "1080", "avc1")
}

// v1.1.0 - with custom directory
func downloadVideo(url string, outputDir string) (string, error) {
    return downloader.DownloadVideoToDir(
        url, "mp4", "1080", "avc1",
        outputDir,
    )
}
```

#### Scenario 3: Adding Progress Tracking

```go
// v1.0.0
func downloadVideo(url string) (string, error) {
    return downloader.DownloadVideo(url, "mp4", "1080", "avc1")
}

// v1.1.0 - with progress
func downloadVideo(url string) (string, error) {
    callback := func(p downloader.DownloadProgress) {
        log.Printf("[%s] %s", time.Now().Format("15:04:05"), p.Stage)
    }
    
    return downloader.DownloadVideoWithProgress(
        url, "mp4", "1080", "avc1",
        callback,
    )
}
```

#### Scenario 4: Adding Metadata Check

```go
// v1.1.0 - check metadata before downloading
func downloadVideoIfSuitable(url string) (string, error) {
    // Fetch metadata first (fast, no download)
    metadata, err := downloader.GetVideoMetadata(url)
    if err != nil {
        return "", err
    }
    
    // Check if video meets criteria
    if metadata.Duration > 3600 { // > 1 hour
        return "", fmt.Errorf("video too long: %s", metadata.DurationString)
    }
    
    if metadata.FilesizeApprox > 500*1024*1024 { // > 500MB
        return "", fmt.Errorf("file too large: %.2f MB", 
            float64(metadata.FilesizeApprox)/(1024*1024))
    }
    
    // Download if suitable
    return downloader.DownloadVideo(url, "mp4", "1080", "avc1")
}
```

### Breaking Changes

**None!** Version 1.1.0 is fully backward compatible.

### Deprecations

**None!** All v1.0.0 functions remain fully supported.

### Performance Improvements

All these improvements happen automatically:

| Aspect | v1.0.0 | v1.1.0 |
|--------|--------|--------|
| Memory usage (large files) | High, can overflow | Constant, streaming |
| Download speed | Standard | Faster with concurrent fragments |
| Reliability | Basic | Automatic retries, better timeouts |
| Progress visibility | None | Optional callbacks |
| Large file support | Limited | Unlimited |
| Binary management | Manual | Automatic |

### Troubleshooting

**"Downloads are slower"** (unlikely):

If you experience slower downloads, you can tune the settings:

```go
// Reduce concurrency if network is unstable
downloader.SetMaxConcurrentDownloads(1)

// Reduce chunk size if memory is limited
downloader.SetChunkSize(16 * 1024 * 1024) // 16MB
```

**"Progress callback not working"**:

Make sure you're using the `*WithProgress` variants:

```go
// Wrong - callback is ignored
downloader.DownloadVideo(url, "", "", "", callback) // ❌

// Correct
downloader.DownloadVideoWithProgress(url, "", "", "", callback) // ✅
```

**"Auto-installation not working"**:

Check if you've disabled it:

```bash
# Check environment variable
echo $GOSTREAMPULLER_NO_AUTO_INSTALL

# If set to 1, unset it or install manually
unset GOSTREAMPULLER_NO_AUTO_INSTALL
gostreampuller-cli setup
```

### Testing Your Upgrade

1. **No code changes**: Just update your `go.mod`:
   ```bash
   go get github.com/debargha2001/gostreampuller@v1.1.0
   ```

2. **Run your existing tests**: They should all pass

3. **Test with a large file**: Try a 30+ minute video to see improvements

4. **Optional**: Add progress tracking or custom directories for better UX

### Summary

**For most users**: Just update and enjoy better performance automatically! 🚀

**For power users**: Explore new features:
- Custom output directories
- Progress tracking
- Metadata fetching
- Performance tuning

**For everyone**: Your existing code continues to work without any changes.

---

## 🎯 Performance Tips for Large Files

### Best Practices

1. **For very large files (1+ hour, 4K)**:
   ```go
   downloader.SetChunkSize(64 * 1024 * 1024) // 64MB
   downloader.SetMaxConcurrentDownloads(5)
   ```

2. **For moderate files (10-30 min)**:
   ```go
   // Use defaults (32MB chunks, 3 concurrent downloads)
   downloader.DownloadVideo(url, "mp4", "1080", "avc1")
   ```

3. **Monitor progress for long downloads**:
   ```go
   downloader.DownloadVideoWithProgress(url, format, resolution, codec, progressCallback)
   ```

### Features Handling Large Files

- **Streaming Architecture**: Files are processed in chunks, not loaded entirely into memory
- **Concurrent Fragments**: Multiple fragments downloaded simultaneously
- **Automatic Timeouts**: 30-minute timeout for downloads, 20-minute for conversions
- **Retry Logic**: 10 retries for both downloads and fragments
- **Optimized FFmpeg**: Uses `-movflags +faststart` and `-max_muxing_queue_size` for large files
- **Buffer Management**: Configurable buffer sizes prevent memory overflow

---

## 📣 License

MIT License © [Debargha Dutta](https://github.com/debarghadutta)
