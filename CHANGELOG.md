# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-11-19

### Added

- **Video Metadata Fetching**: New `GetVideoMetadata()` function to fetch comprehensive video information without downloading
  - Retrieve title, description, uploader, duration, views, likes, and more
  - Technical details: resolution, codecs, file size, FPS
  - Channel information and timestamps
  - Live stream detection (`IsLive`, `WasLive`, `LiveStatus`)
  - Categories, tags, and subtitles
  - Raw metadata access for custom fields
  - Context-aware variant: `GetVideoMetadataWithContext()`
  - Comprehensive `VideoMetadata` struct with 40+ fields
- **Metadata Tests**: Complete test suite for metadata functionality in `test/metadata_test.go`
  - Basic metadata fetching
  - Context and timeout handling
  - Invalid URL validation
  - Multiple platform support
  - Live stream detection
  - Field validation tests
- **Metadata Examples**: Practical examples in `test/examples/metadata/main.go`
  - Fetch video information before downloading
  - Check quality and size before download
  - Smart download decisions based on metadata
  - Access raw metadata fields
- **CLI Setup Tool**: New `gostreampuller-cli` command-line tool for automated dependency installation
  - `gostreampuller-cli setup`: Automatically downloads and installs yt-dlp and ffmpeg
  - `gostreampuller-cli check`: Verifies dependency installation status
  - Cross-platform support (Linux, macOS, Windows)
  - Binaries installed to `~/.gostreampuller/bin/`
- **Automated Binary Detection**: Package automatically detects and uses locally installed binaries
  - No manual PATH configuration needed
  - Falls back to system PATH if local binaries not found
  - Perfect for CI/CD and offline environments
- **Production-Grade Test Suite**: Comprehensive testing infrastructure in `test/` directory
  - `test/integration_test.go`: Integration tests with real downloads
  - `test/unit_test.go`: Fast unit tests
  - `test/run_tests.sh`: Automated test runner with multiple options
  - `test/examples/`: Practical usage examples
  - `test/README.md`: Complete testing documentation
- **Internal Installer Package**: `internal/installer` for cross-platform binary installation
  - Downloads latest stable releases of yt-dlp and ffmpeg
  - Handles platform-specific installation (Linux, macOS, Windows)
  - Progress reporting during downloads
  - Automatic extraction from archives
- **Large File Support**: Complete rewrite to handle videos of any size (30+ minutes, 4K, multi-hour content)
- **Streaming Architecture**: Implemented streaming and chunked processing to prevent memory overflow
- **Concurrent Downloads**: Added support for downloading video fragments concurrently (configurable, default: 3)
- **Progress Tracking**: New `DownloadVideoWithProgress()` and `DownloadAudioWithProgress()` functions with callback support
- **Progress Callbacks**: `ProgressCallback` type to monitor download and conversion stages
- **Timeout Protection**: Automatic timeout handling (30 min for downloads, 20 min for conversions)
- **Retry Logic**: Built-in retry mechanism (10 retries for downloads and fragments)
- **Configurable Settings**:
  - `SetChunkSize()`: Customize buffer size for streaming operations (default: 32MB)
  - `SetMaxConcurrentDownloads()`: Control concurrent fragment downloads (default: 3)
- **Optimized yt-dlp flags**:
  - `--no-part`: Avoid .part files for large downloads
  - `--concurrent-fragments`: Download fragments in parallel
  - `--buffer-size`: Optimized buffer management
  - `--retries` and `--fragment-retries`: Automatic retry on failure
- **Optimized FFmpeg flags**:
  - `-movflags +faststart`: Optimize for streaming
  - `-max_muxing_queue_size`: Handle large file muxing
- **Helper Functions**:
  - `streamCommand()`: Execute commands with streaming output handling
  - `copyFileStreaming()`: Copy large files efficiently using buffered I/O

### Changed

- **Project Structure**: Reorganized for production use
  - Moved all tests to `test/` directory
  - Moved examples to `test/examples/`
  - Added `cmd/` for CLI tools
  - Added `internal/` for internal packages
- `DownloadVideo()`: Now uses streaming architecture internally (backward compatible)
- `DownloadAudio()`: Now uses streaming architecture internally (backward compatible)
- Audio download now detects actual downloaded format instead of assuming webm
- Improved error messages with more context
- Binary paths now auto-detected on package initialization

### Improved

- Memory efficiency: Large files no longer cause memory overflow
- Performance: Concurrent fragment downloads significantly speed up large file downloads
- Reliability: Retry logic handles temporary network issues automatically
- User experience: Progress callbacks provide visibility into long-running operations
- Developer experience: One-command setup for all dependencies
- Testing: Comprehensive test suite with automated runner

### Technical Details

- Added `context.Context` support for timeout management
- Implemented `sync.WaitGroup` for goroutine coordination
- Added `bufio.Scanner` with configurable buffer sizes for streaming
- Used `io.CopyBuffer` for efficient large file copying
- Cross-platform binary installation with proper error handling

## [1.0.0] - Initial Release

### Added

- Basic video download functionality
- Basic audio extraction functionality
- Support for multiple video formats (mp4, webm, mkv, avi, mov, flv)
- Support for multiple audio formats (mp3, aac, ogg, wav, flac, m4a)
- Customizable video resolution (144p to 4K)
- Customizable video codec (H.264, VP9, VP8, AV1, H.265)
- Customizable audio codec and bitrate
- Cross-platform support (Linux, macOS, Windows)
- Binary path overrides for yt-dlp and ffmpeg
- Automatic temporary file cleanup
- Smart format detection and conversion
