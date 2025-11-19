#!/bin/bash

# gostreampuller Test Suite Runner
# This script runs all tests and examples for the gostreampuller package

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}╔════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   gostreampuller Test Suite Runner        ║${NC}"
echo -e "${BLUE}╔════════════════════════════════════════════╗${NC}"
echo ""

# Function to print section header
print_section() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

# Function to check command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_section "🔍 Checking Prerequisites"

if ! command_exists go; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Go is installed: $(go version)${NC}"

# Check if binaries are available
if command_exists yt-dlp || [ -f "$HOME/.gostreampuller/bin/yt-dlp" ]; then
    echo -e "${GREEN}✓ yt-dlp is available${NC}"
else
    echo -e "${YELLOW}⚠️  yt-dlp not found. Run 'gostreampuller-cli setup' or install manually${NC}"
fi

if command_exists ffmpeg || [ -f "$HOME/.gostreampuller/bin/ffmpeg" ]; then
    echo -e "${GREEN}✓ ffmpeg is available${NC}"
else
    echo -e "${YELLOW}⚠️  ffmpeg not found. Run 'gostreampuller-cli setup' or install manually${NC}"
fi

# Change to root directory
cd "$ROOT_DIR"

# Run build check
print_section "🔨 Building Package"
echo "Building all packages..."
if go build ./...; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Run unit tests
print_section "🧪 Running Unit Tests"
echo "Running unit tests (fast, no downloads)..."
if go test -v -short ./test -run "TestPackage|TestSetters|TestDownloadProgress"; then
    echo -e "${GREEN}✓ Unit tests passed${NC}"
else
    echo -e "${RED}❌ Unit tests failed${NC}"
    exit 1
fi

# Run integration tests (only if -integration flag is passed)
if [[ "$*" == *"--integration"* ]] || [[ "$*" == *"-i"* ]]; then
    print_section "🌐 Running Integration Tests"
    echo "Running integration tests (requires internet and binaries)..."
    echo "This will download actual videos and may take several minutes..."
    echo ""
    
    if go test -v ./test -timeout 30m; then
        echo -e "${GREEN}✓ Integration tests passed${NC}"
    else
        echo -e "${RED}❌ Integration tests failed${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}ℹ️  Skipping integration tests (use --integration flag to run them)${NC}"
fi

# Run examples (only if -examples flag is passed)
if [[ "$*" == *"--examples"* ]] || [[ "$*" == *"-e"* ]]; then
    print_section "💡 Running Examples"
    echo "Building and running example programs..."
    echo ""
    
    if [ -d "$SCRIPT_DIR/examples" ]; then
        cd "$SCRIPT_DIR/examples"
        if go run main.go; then
            echo -e "${GREEN}✓ Examples completed successfully${NC}"
        else
            echo -e "${RED}❌ Examples failed${NC}"
            exit 1
        fi
        cd "$ROOT_DIR"
    else
        echo -e "${YELLOW}⚠️  No examples directory found${NC}"
    fi
else
    echo -e "${YELLOW}ℹ️  Skipping examples (use --examples flag to run them)${NC}"
fi

# Run benchmarks (only if -bench flag is passed)
if [[ "$*" == *"--bench"* ]] || [[ "$*" == *"-b"* ]]; then
    print_section "⚡ Running Benchmarks"
    echo "Running performance benchmarks..."
    echo ""
    
    go test -bench=. -benchmem ./test
    echo -e "${GREEN}✓ Benchmarks completed${NC}"
else
    echo -e "${YELLOW}ℹ️  Skipping benchmarks (use --bench flag to run them)${NC}"
fi

# Run code coverage (only if -coverage flag is passed)
if [[ "$*" == *"--coverage"* ]] || [[ "$*" == *"-c"* ]]; then
    print_section "📊 Generating Code Coverage"
    echo "Generating coverage report..."
    
    go test -coverprofile=coverage.out ./downloader ./test
    go tool cover -html=coverage.out -o coverage.html
    
    echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
    
    # Show coverage summary
    echo ""
    go tool cover -func=coverage.out | tail -1
else
    echo -e "${YELLOW}ℹ️  Skipping coverage (use --coverage flag to generate)${NC}"
fi

# Final summary
print_section "📝 Test Summary"
echo -e "${GREEN}✅ All tests completed successfully!${NC}"
echo ""
echo "Quick reference:"
echo "  • Unit tests: fast, no downloads required"
echo "  • Integration tests: requires internet and binaries (--integration)"
echo "  • Examples: runs example programs (--examples)"
echo "  • Benchmarks: performance testing (--bench)"
echo "  • Coverage: code coverage report (--coverage)"
echo ""
echo "Usage examples:"
echo "  ./run_tests.sh                    # Run basic tests only"
echo "  ./run_tests.sh --integration      # Include integration tests"
echo "  ./run_tests.sh -i -e              # Integration tests + examples"
echo "  ./run_tests.sh -i -e -b -c        # Everything + benchmarks + coverage"
echo ""
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}🎉 Test suite execution complete!${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
