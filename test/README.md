# Test Suite for gostreampuller

This directory contains the comprehensive test suite for the gostreampuller package.

## 📁 Structure

```
test/
├── run_tests.sh          # Main test runner script
├── integration_test.go   # Integration tests (requires internet & binaries)
├── unit_test.go          # Unit tests (fast, no dependencies)
├── metadata_test.go      # Metadata fetching tests
└── examples/
    ├── main.go           # Basic usage demonstrations
    └── metadata/
        └── main.go       # Metadata fetching examples
```

## 🚀 Quick Start

### Run Basic Tests

```bash
cd test
./run_tests.sh
```

This runs:
- ✅ Build verification
- ✅ Unit tests (fast, no downloads)

### Run All Tests

```bash
./run_tests.sh --integration --examples
```

Or use short flags:
```bash
./run_tests.sh -i -e
```

## 📋 Test Commands

### Unit Tests Only (Fast)
```bash
go test -v -short ./test
```

### Integration Tests (Requires Internet)
```bash
go test -v ./test -timeout 30m
```

### Specific Test
```bash
go test -v ./test -run TestDownloadVideo
```

### With Coverage
```bash
go test -coverprofile=coverage.out ./test
go tool cover -html=coverage.out
```

### Benchmarks
```bash
go test -bench=. -benchmem ./test
```

## 🎯 Test Flags

The `run_tests.sh` script supports these flags:

| Flag | Short | Description |
|------|-------|-------------|
| `--integration` | `-i` | Run integration tests (downloads actual videos) |
| `--examples` | `-e` | Run example programs |
| `--bench` | `-b` | Run performance benchmarks |
| `--coverage` | `-c` | Generate code coverage report |

### Examples

```bash
# Run everything
./run_tests.sh -i -e -b -c

# Just integration tests
./run_tests.sh --integration

# Integration + examples
./run_tests.sh -i -e

# Quick unit tests only
./run_tests.sh
```

## 📊 Test Categories

### 1. Unit Tests (`unit_test.go`)

Fast tests that don't require downloads:
- Package initialization
- Configuration setters
- Data structure validation

**Run with:**
```bash
go test -v -short ./test
```

### 2. Integration Tests (`integration_test.go`)

Real download tests (requires internet and binaries):
- Video downloads
- Audio downloads
- Progress tracking
- Concurrent downloads
- Error handling

**Run with:**
```bash
go test -v ./test
```

### 3. Metadata Tests (`metadata_test.go`)

Video metadata fetching tests (requires internet):
- Basic metadata fetching
- Metadata with context/timeout
- Invalid URL handling
- Multiple platform support
- Live stream detection
- Metadata field validation

**Run with:**
```bash
go test -v -run TestGetVideoMetadata ./test
```

### 4. Examples

Practical usage demonstrations:

**Basic Usage** (`examples/main.go`):
- Basic video download
- Progress tracking
- Audio extraction
- Large file optimization

**Metadata Example** (`examples/metadata/main.go`):
- Fetch video information
- Check quality before download
- Detect live streams
- Smart download decisions
- Access raw metadata

**Run with:**
```bash
# Basic examples
cd examples && go run main.go

# Metadata examples
cd examples/metadata && go run main.go
```

## ⚙️ Prerequisites

### Required
- Go 1.24.1 or later
- Internet connection (for integration tests)

### Dependencies (for integration tests)
- `yt-dlp` - Install with: `gostreampuller-cli setup`
- `ffmpeg` - Install with: `gostreampuller-cli setup`

Or install manually:
```bash
# macOS
brew install yt-dlp ffmpeg

# Linux
sudo apt install yt-dlp ffmpeg

# Windows
choco install yt-dlp ffmpeg
```

## 🔧 CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      # Install dependencies
      - name: Install yt-dlp and ffmpeg
        run: |
          sudo apt-get update
          sudo apt-get install -y yt-dlp ffmpeg
      
      # Run tests
      - name: Run unit tests
        run: cd test && ./run_tests.sh
      
      - name: Run integration tests
        run: cd test && ./run_tests.sh --integration
```

## 📝 Writing New Tests

### Unit Test Template

```go
func TestNewFeature(t *testing.T) {
    // Setup
    originalValue := downloader.SomeSetting
    defer func() { downloader.SomeSetting = originalValue }()
    
    // Test
    result := downloader.NewFeature()
    
    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
    
    t.Log("✓ New feature works correctly")
}
```

### Integration Test Template

```go
func TestNewDownloadFeature(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Test download
    file, err := downloader.NewDownloadFeature(url)
    if err != nil {
        t.Fatalf("Failed: %v", err)
    }
    defer os.Remove(file)
    
    // Verify
    if _, err := os.Stat(file); os.IsNotExist(err) {
        t.Fatal("File not created")
    }
    
    t.Logf("✓ Successfully tested new feature: %s", file)
}
```

## 🐛 Debugging Failed Tests

### Test fails with "yt-dlp not found"
```bash
# Check if installed locally
ls ~/.gostreampuller/bin/

# Run setup
go run cmd/gostreampuller-cli/main.go setup

# Or install system-wide
brew install yt-dlp  # macOS
sudo apt install yt-dlp  # Linux
```

### Test fails with network error
```bash
# Check internet connection
ping youtube.com

# Try with longer timeout
go test -v ./test -timeout 60m
```

### Test downloads fail
```bash
# Update yt-dlp to latest
yt-dlp -U

# Or reinstall
go run cmd/gostreampuller-cli/main.go setup
```

## 📈 Performance Benchmarks

Run benchmarks to measure performance:

```bash
# Run all benchmarks
go test -bench=. ./test

# Run specific benchmark
go test -bench=BenchmarkDownloadSmallVideo ./test

# With memory profiling
go test -bench=. -benchmem ./test

# Save results
go test -bench=. ./test > bench.txt
```

## 🎯 Test Coverage Goals

| Component | Target Coverage |
|-----------|----------------|
| downloader package | > 80% |
| installer package | > 70% |
| CLI commands | > 60% |

Check current coverage:
```bash
./run_tests.sh --coverage
open coverage.html
```

## 📞 Support

If tests fail:
1. Check prerequisites are installed
2. Verify internet connection
3. Update yt-dlp: `yt-dlp -U`
4. Check GitHub issues
5. Run with verbose flag: `go test -v ./test`

## 🔄 Continuous Testing

For development, use watch mode:

```bash
# Install watchexec
brew install watchexec

# Auto-run tests on file changes
watchexec -w downloader -w test "cd test && ./run_tests.sh"
```

---

**Note**: Integration tests download actual videos and may take several minutes. Use `-short` flag or skip `--integration` for faster testing during development.
