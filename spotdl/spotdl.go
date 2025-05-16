// Package spotdl provides functionality to download Spotify tracks
package spotdl

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// filenameFormat defines the output filename format
	filenameFormat = "{artist} - {title}.{output-ext}"
)

//go:embed assets/spotdl_linux assets/windows/spotdl.exe assets/windows/ffmpeg.exe assets/windows/ffplay.exe assets/windows/ffprobe.exe
var spotdlBinary embed.FS

// Downloader handles downloading of tracks from Spotify
type Downloader struct {
	ctx context.Context
}

// NewDownloader creates a new Downloader instance
func NewDownloader() *Downloader {
	return &Downloader{}
}

// Startup is called when the application starts
func (d *Downloader) Startup(ctx context.Context) {
	d.ctx = ctx
}

// Download downloads a track from the provided Spotify link
// Parameters:
// - link: Spotify link to download from
// - outputPath: directory where to save the downloaded files
// - format: output format (mp3, wav, etc.)
// - bitrate: quality of the output (128k, 320k, etc.)
// - songsToDelete: optional list of songs to delete after download
// Returns: boolean indicating whether the download was successful
func (d *Downloader) Download(link, outputPath, format, bitrate string, songsToDelete []string) bool {
	// Extract the spotdl binary to a temporary location
	tmpDir, err := os.MkdirTemp("", "spotdl")
	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to create temp directory: %v", err))
		return false
	}
	defer os.RemoveAll(tmpDir)

	// Determine file paths based on OS
	var spotdlPath string
	var binPath string
	var ffmpegPath, ffplayPath, ffprobePath string
	isWindows := runtime.GOOS == "windows"

	if isWindows {
		binPath = "assets/windows/spotdl.exe"
		spotdlPath = filepath.Join(tmpDir, "spotdl.exe")

		// Extract FFmpeg binaries for Windows
		ffmpegPath = filepath.Join(tmpDir, "ffmpeg.exe")
		ffplayPath = filepath.Join(tmpDir, "ffplay.exe")
		ffprobePath = filepath.Join(tmpDir, "ffprobe.exe")

		// Extract FFmpeg binaries
		if err := d.extractBinary("assets/windows/ffmpeg.exe", ffmpegPath); err != nil {
			d.emitErrorEvent(fmt.Sprintf("failed to extract ffmpeg binary: %v", err))
			return false
		}
		if err := d.extractBinary("assets/windows/ffplay.exe", ffplayPath); err != nil {
			d.emitErrorEvent(fmt.Sprintf("failed to extract ffplay binary: %v", err))
			return false
		}
		if err := d.extractBinary("assets/windows/ffprobe.exe", ffprobePath); err != nil {
			d.emitErrorEvent(fmt.Sprintf("failed to extract ffprobe binary: %v", err))
			return false
		}
	} else {
		binPath = "assets/spotdl_linux"
		spotdlPath = filepath.Join(tmpDir, "spotdl")
	}

	// Get the embedded binary
	if err := d.extractBinary(binPath, spotdlPath); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to extract spotdl binary: %v", err))
		return false
	}

	// Prepare arguments
	args := []string{
		link,
		"--bitrate", bitrate,
		"--format", format,
		"--print-errors",
		"--output",
	}

	outputFilePath := filenameFormat
	if outputPath != "" {
		outputFilePath = filepath.Join(outputPath, filenameFormat)
	}
	args = append(args, outputFilePath)

	// Execute spotdl
	cmd := exec.Command(spotdlPath, args...)

	// Set up stdout and stderr pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to create stdout pipe: %v", err))
		return false
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to create stderr pipe: %v", err))
		return false
	}

	// Start the command
	wailsruntime.EventsEmit(d.ctx, "update_in_download", "Downloading")
	if err := cmd.Start(); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to start command: %v", err))
		return false
	}

	// Create a wait group to wait for the goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Process stdout
	go d.pipeReader(&wg, stdoutPipe)

	// Process stderr
	go d.pipeReader(&wg, stderrPipe)

	// Wait for the command to finish
	err = cmd.Wait()
	wg.Wait()

	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("command execution failed: %v", err))
		return false
	}

	wailsruntime.EventsEmit(d.ctx, "update_in_download", "Done")
	return true
}

// extractBinary extracts a binary from the embedded FS to the target path
func (d *Downloader) extractBinary(embeddedPath, targetPath string) error {
	binData, err := spotdlBinary.ReadFile(embeddedPath)
	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to read embedded binary %s: %v", embeddedPath, err))
		return fmt.Errorf("failed to read embedded binary %s: %w", embeddedPath, err)
	}

	// Write the binary to the temp location
	if err := os.WriteFile(targetPath, binData, 0755); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to write binary to temp location: %v", err))
		return fmt.Errorf("failed to write binary to temp location: %w", err)
	}

	return nil
}

// pipeReader reads from a pipe and emits events with the content
func (d *Downloader) pipeReader(wg *sync.WaitGroup, pipe io.ReadCloser) {
	defer wg.Done()

	buf := make([]byte, 1024)
	for {
		n, err := pipe.Read(buf)
		if n > 0 {
			output := string(buf[:n])
			log.Print(output)
			wailsruntime.EventsEmit(d.ctx, "update_in_download", output)
		}
		if err != nil {
			break
		}
	}
}

// emitErrorEvent emits a fatal error event and a done event
func (d *Downloader) emitErrorEvent(errMsg string) {
	wailsruntime.EventsEmit(d.ctx, "update_in_download", fmt.Sprintf("fatal_error: %s", errMsg))
	wailsruntime.EventsEmit(d.ctx, "update_in_download", "Done")
}
