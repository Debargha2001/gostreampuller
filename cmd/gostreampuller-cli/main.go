package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/debargha2001/gostreampuller/internal/installer"
)

const version = "1.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "setup":
		handleSetup()
	case "check":
		handleCheck()
	case "version", "--version", "-v":
		fmt.Printf("gostreampuller-cli v%s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`gostreampuller-cli - Setup and manage gostreampuller dependencies

USAGE:
    gostreampuller-cli <command>

COMMANDS:
    setup      Install yt-dlp and ffmpeg locally
    check      Check if dependencies are installed
    version    Show version information
    help       Show this help message

EXAMPLES:
    # Install all dependencies
    gostreampuller-cli setup

    # Check installation status
    gostreampuller-cli check

DEPENDENCIES:
    The setup command will download and install:
    - yt-dlp: Latest stable release
    - ffmpeg: Latest stable release

    Binaries are installed to: ~/.gostreampuller/bin/

NOTES:
    - On macOS, you may want to use Homebrew: brew install yt-dlp ffmpeg
    - On Linux, package managers also work: apt/yum install yt-dlp ffmpeg
    - Windows users can use Chocolatey or Scoop, or this setup command
    - The setup command works on all platforms for offline installations`)
}

func handleSetup() {
	fmt.Println("🚀 gostreampuller Setup")
	fmt.Println("=" + "=============================")
	fmt.Println()

	// Check current installation status
	ytdlpInstalled, ffmpegInstalled, err := installer.CheckInstallation()
	if err != nil {
		fmt.Printf("Error checking installation: %v\n", err)
	}

	binDir, err := installer.GetBinariesDir()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Installation directory: %s\n\n", binDir)

	// Install yt-dlp
	if ytdlpInstalled {
		fmt.Println("✓ yt-dlp is already installed")
	} else {
		fmt.Println("📥 Installing yt-dlp...")
		progressFn := func(msg string) {
			fmt.Printf("   %s\n", msg)
		}
		
		if err := installer.InstallYTDLP(progressFn); err != nil {
			fmt.Printf("❌ Failed to install yt-dlp: %v\n", err)
			fmt.Println("\nTroubleshooting:")
			fmt.Println("  - Check your internet connection")
			fmt.Println("  - Try installing manually from: https://github.com/yt-dlp/yt-dlp")
			os.Exit(1)
		}
		fmt.Println()
	}

	// Install ffmpeg
	if ffmpegInstalled {
		fmt.Println("✓ ffmpeg is already installed")
	} else {
		fmt.Println("📥 Installing ffmpeg...")
		progressFn := func(msg string) {
			fmt.Printf("   %s\n", msg)
		}
		
		if err := installer.InstallFFMPEG(progressFn); err != nil {
			fmt.Printf("❌ Failed to install ffmpeg: %v\n", err)
			fmt.Println("\nTroubleshooting:")
			fmt.Println("  - Check your internet connection")
			fmt.Println("  - macOS users: Try 'brew install ffmpeg'")
			fmt.Println("  - Linux users: Try 'sudo apt install ffmpeg' or 'sudo yum install ffmpeg'")
			fmt.Println("  - Windows users: Try 'choco install ffmpeg' or download from https://ffmpeg.org")
			os.Exit(1)
		}
		fmt.Println()
	}

	// Verify installation
	fmt.Println("🔍 Verifying installation...")
	ytdlpPath, _ := installer.GetYTDLPPath()
	ffmpegPath, _ := installer.GetFFMPEGPath()

	fmt.Printf("   yt-dlp: %s\n", ytdlpPath)
	fmt.Printf("   ffmpeg: %s\n", ffmpegPath)
	fmt.Println()
	
	// Create marker file to indicate CLI installation was done
	// This prevents auto-installation from triggering later
	createCLIMarker()

	fmt.Println("✅ Setup complete!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. The gostreampuller package will automatically use these binaries")
	fmt.Println("  2. No need to set PATH or configure anything")
	fmt.Println("  3. Start downloading videos with your Go code!")
	fmt.Println()
	fmt.Println("Example usage:")
	fmt.Println(`  import "github.com/debargha2001/gostreampuller/downloader"`)
	fmt.Println(`  video, err := downloader.DownloadVideo(url, "mp4", "1080", "avc1")`)
}

// createCLIMarker creates a marker file to indicate CLI setup was used
func createCLIMarker() {
	binDir, err := installer.GetBinariesDir()
	if err != nil {
		return
	}
	
	configDir := filepath.Dir(binDir)
	markerFile := filepath.Join(configDir, ".cli_installed")
	
	// Create marker file with timestamp
	timestamp := []byte(fmt.Sprintf("cli-installed: %s", time.Now().Format(time.RFC3339)))
	os.WriteFile(markerFile, timestamp, 0644)
}

func handleCheck() {
	fmt.Println("🔍 Checking gostreampuller dependencies...")
	fmt.Println()

	ytdlpInstalled, ffmpegInstalled, err := installer.CheckInstallation()
	if err != nil {
		fmt.Printf("Error checking installation: %v\n", err)
		os.Exit(1)
	}

	binDir, _ := installer.GetBinariesDir()
	fmt.Printf("Binary directory: %s\n\n", binDir)

	// Check yt-dlp
	if ytdlpInstalled {
		ytdlpPath, _ := installer.GetYTDLPPath()
		fmt.Printf("✅ yt-dlp: Installed\n")
		fmt.Printf("   Path: %s\n", ytdlpPath)
	} else {
		fmt.Printf("❌ yt-dlp: Not installed\n")
		fmt.Printf("   Run: gostreampuller-cli setup\n")
	}
	fmt.Println()

	// Check ffmpeg
	if ffmpegInstalled {
		ffmpegPath, _ := installer.GetFFMPEGPath()
		fmt.Printf("✅ ffmpeg: Installed\n")
		fmt.Printf("   Path: %s\n", ffmpegPath)
	} else {
		fmt.Printf("❌ ffmpeg: Not installed\n")
		fmt.Printf("   Run: gostreampuller-cli setup\n")
	}
	fmt.Println()

	// Overall status
	if ytdlpInstalled && ffmpegInstalled {
		fmt.Println("✅ All dependencies are installed and ready to use!")
	} else {
		fmt.Println("⚠️  Some dependencies are missing. Run 'gostreampuller-cli setup' to install them.")
		os.Exit(1)
	}
}
