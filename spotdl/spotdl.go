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
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// filenameFormat defines the output filename format
	filenameFormat = "{artist} - {title}.{output-ext}"
)

//go:embed assets/spotdl_linux
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
func (d *Downloader) Download(link, outputPath, format, bitrate string, songsToDelete []string) error {
	// Extract the spotdl binary to a temporary location
	tmpDir, err := os.MkdirTemp("", "spotdl")
	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to create temp directory: %v", err))
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	spotdlPath := filepath.Join(tmpDir, "spotdl")

	// Get the embedded binary
	binData, err := spotdlBinary.ReadFile("assets/spotdl_linux")
	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to read embedded binary: %v", err))
		return fmt.Errorf("failed to read embedded binary: %w", err)
	}

	// Write the binary to the temp location
	if err := os.WriteFile(spotdlPath, binData, 0755); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to write binary to temp location: %v", err))
		return fmt.Errorf("failed to write binary to temp location: %w", err)
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
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to create stderr pipe: %v", err))
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	runtime.EventsEmit(d.ctx, "update_in_download", "Downloading")
	if err := cmd.Start(); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to start command: %v", err))
		return fmt.Errorf("failed to start command: %w", err)
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
		return fmt.Errorf("command execution failed: %w", err)
	}

	runtime.EventsEmit(d.ctx, "update_in_download", "Done")
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
			runtime.EventsEmit(d.ctx, "update_in_download", output)
		}
		if err != nil {
			break
		}
	}
}

// emitErrorEvent emits a fatal error event and a done event
func (d *Downloader) emitErrorEvent(errMsg string) {
	runtime.EventsEmit(d.ctx, "update_in_download", fmt.Sprintf("fatal_error: %s", errMsg))
	runtime.EventsEmit(d.ctx, "update_in_download", "Done")
}
