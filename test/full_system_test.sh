#!/bin/bash

# Full System Test for gostreampuller v1.1.0
# Tests: Binary installation, Metadata fetching, Video/Audio downloads

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_SKIPPED=0

# Helper functions
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_skip() {
    echo -e "${YELLOW}⊘ $1${NC}"
    ((TESTS_SKIPPED++))
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Change to project root
cd "$(dirname "$0")/.."

print_header "Full System Test Suite"
print_info "Testing: Binary installation, Metadata, Downloads"
print_info "Project: gostreampuller v1.1.0"
print_info "Date: $(date)"

# Test 1: Build the CLI tool
print_header "1. Building CLI Tool"
if go build -o /tmp/gostreampuller-cli ./cmd/gostreampuller-cli; then
    print_success "CLI tool built successfully"
else
    print_error "Failed to build CLI tool"
    exit 1
fi

# Test 2: Check CLI version
print_header "2. Testing CLI Version Command"
if /tmp/gostreampuller-cli version; then
    print_success "CLI version command works"
else
    print_error "CLI version command failed"
fi

# Test 3: Check current installation status
print_header "3. Checking Current Installation Status"
if /tmp/gostreampuller-cli check; then
    print_success "Dependencies check completed"
else
    print_warning "Some dependencies may be missing (this is OK)"
fi

# Test 4: Install binaries (if user confirms)
print_header "4. Binary Installation Test"
echo -e "${YELLOW}Do you want to test binary installation? This will install yt-dlp and ffmpeg to ~/.gostreampuller/bin/${NC}"
echo -e "${YELLOW}If already installed, it will verify them. (y/N)${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    print_info "Installing binaries..."
    if /tmp/gostreampuller-cli setup; then
        print_success "Binary installation completed"
        
        # Verify binaries exist
        if [ -f "$HOME/.gostreampuller/bin/yt-dlp" ] || [ -f "$HOME/.gostreampuller/bin/yt-dlp.exe" ]; then
            print_success "yt-dlp binary found in ~/.gostreampuller/bin/"
        else
            print_error "yt-dlp binary not found after installation"
        fi
        
        if [ -f "$HOME/.gostreampuller/bin/ffmpeg" ] || [ -f "$HOME/.gostreampuller/bin/ffmpeg.exe" ]; then
            print_success "ffmpeg binary found in ~/.gostreampuller/bin/"
        else
            print_error "ffmpeg binary not found after installation"
        fi
    else
        print_error "Binary installation failed"
    fi
else
    print_skip "Binary installation test skipped"
fi

# Test 5: Build verification
print_header "5. Building All Packages"
if go build ./...; then
    print_success "All packages built successfully"
else
    print_error "Package build failed"
    exit 1
fi

# Test 6: Unit tests
print_header "6. Running Unit Tests"
if go test -v -short ./test -timeout 2m; then
    print_success "Unit tests passed"
else
    print_error "Unit tests failed"
fi

# Test 7: Metadata tests (requires internet)
print_header "7. Running Metadata Tests"
echo -e "${YELLOW}Do you want to run metadata tests? (requires internet) (y/N)${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    print_info "Running metadata tests..."
    if go test -v -run TestGetVideoMetadata ./test -timeout 5m; then
        print_success "Metadata tests passed"
    else
        print_error "Metadata tests failed"
    fi
else
    print_skip "Metadata tests skipped"
fi

# Test 8: Integration tests (requires internet and binaries)
print_header "8. Running Integration Tests"
echo -e "${YELLOW}Do you want to run integration tests? (requires internet, binaries, and downloads ~50MB) (y/N)${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    print_info "Running integration tests (this may take several minutes)..."
    if go test -v -run "TestDownload" ./test -timeout 30m; then
        print_success "Integration tests passed"
    else
        print_error "Integration tests failed"
    fi
    
    # Clean up downloaded files
    print_info "Cleaning up test downloads..."
    rm -f *.mp4 *.mp3 *.webm *.mkv 2>/dev/null || true
    print_success "Cleanup completed"
else
    print_skip "Integration tests skipped"
fi

# Test 9: Run basic example
print_header "9. Testing Basic Example"
echo -e "${YELLOW}Do you want to run the basic example? (downloads a real video) (y/N)${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    print_info "Running basic example..."
    cd test/examples
    if timeout 5m go run main.go; then
        print_success "Basic example executed successfully"
    else
        print_error "Basic example failed"
    fi
    cd ../..
else
    print_skip "Basic example skipped"
fi

# Test 10: Run metadata example
print_header "10. Testing Metadata Example"
echo -e "${YELLOW}Do you want to run the metadata example? (fetches metadata only, no download) (y/N)${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    print_info "Running metadata example..."
    cd test/examples/metadata
    if timeout 2m go run main.go; then
        print_success "Metadata example executed successfully"
    else
        print_error "Metadata example failed"
    fi
    cd ../../..
else
    print_skip "Metadata example skipped"
fi

# Test 11: Quick functionality test
print_header "11. Quick Functionality Test"
echo -e "${YELLOW}Do you want to run a quick live test? (fetches metadata + downloads a 19-second video) (y/N)${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    print_info "Creating quick test..."
    
    # Create a quick test program
    cat > /tmp/gostreampuller_quick_test.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
    url := "https://www.youtube.com/watch?v=jNQXAC9IVRw" // "Me at the zoo" - 19 seconds
    
    fmt.Println("=== Testing Metadata Fetch ===")
    metadata, err := downloader.GetVideoMetadata(url)
    if err != nil {
        log.Fatalf("Metadata fetch failed: %v", err)
    }
    
    fmt.Printf("✓ Title: %s\n", metadata.Title)
    fmt.Printf("✓ Duration: %s\n", metadata.DurationString)
    fmt.Printf("✓ Uploader: %s\n", metadata.Uploader)
    fmt.Printf("✓ Resolution: %dx%d\n", metadata.Width, metadata.Height)
    
    fmt.Println("\n=== Testing Video Download ===")
    progress := func(p downloader.DownloadProgress) {
        fmt.Printf("  Progress: %s\n", p.Stage)
    }
    
    video, err := downloader.DownloadVideoWithProgress(url, "mp4", "360", "avc1", progress)
    if err != nil {
        log.Fatalf("Video download failed: %v", err)
    }
    
    fmt.Printf("✓ Downloaded: %s\n", video)
    
    // Verify file exists
    if _, err := os.Stat(video); err == nil {
        info, _ := os.Stat(video)
        fmt.Printf("✓ File size: %.2f MB\n", float64(info.Size())/(1024*1024))
        
        // Cleanup
        os.Remove(video)
        fmt.Println("✓ Cleanup completed")
    } else {
        log.Fatalf("Downloaded file not found: %v", err)
    }
    
    fmt.Println("\n✓ All quick tests passed!")
}
EOF
    
    print_info "Running quick test (this will download ~2MB)..."
    if cd /tmp && go mod init quicktest 2>/dev/null && go mod edit -replace github.com/debargha2001/gostreampuller="$OLDPWD" && go mod tidy && timeout 5m go run gostreampuller_quick_test.go; then
        print_success "Quick functionality test passed"
    else
        print_error "Quick functionality test failed"
    fi
    
    # Cleanup
    rm -rf /tmp/gostreampuller_quick_test.go /tmp/go.mod /tmp/go.sum 2>/dev/null || true
    cd "$OLDPWD"
else
    print_skip "Quick functionality test skipped"
fi

# Test 12: Binary path override test
print_header "12. Testing Binary Path Override"
print_info "Creating binary path test..."

cat > /tmp/path_test.go << 'EOF'
package main

import (
    "fmt"
    "github.com/debargha2001/gostreampuller/downloader"
)

func main() {
    // Test setting custom paths
    downloader.SetYTDLPPath("/custom/path/yt-dlp")
    downloader.SetFFMPEGPath("/custom/path/ffmpeg")
    fmt.Println("✓ Custom paths set successfully")
    
    // Test reset
    downloader.ResetBinaryPaths()
    fmt.Println("✓ Paths reset to auto-detected values")
    
    fmt.Println("✓ Binary path override test passed")
}
EOF

if cd /tmp && go mod init pathtest 2>/dev/null && go mod edit -replace github.com/debargha2001/gostreampuller="$OLDPWD" && go mod tidy && go run path_test.go 2>/dev/null; then
    print_success "Binary path override test passed"
else
    print_error "Binary path override test failed"
fi

rm -rf /tmp/path_test.go /tmp/go.mod /tmp/go.sum 2>/dev/null || true
cd "$OLDPWD"

# Final summary
print_header "Test Summary"
echo -e "${GREEN}Passed:  $TESTS_PASSED${NC}"
echo -e "${RED}Failed:  $TESTS_FAILED${NC}"
echo -e "${YELLOW}Skipped: $TESTS_SKIPPED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}✓ ALL TESTS PASSED!${NC}"
    echo -e "${GREEN}========================================${NC}"
    exit 0
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}✗ SOME TESTS FAILED${NC}"
    echo -e "${RED}========================================${NC}"
    exit 1
fi
