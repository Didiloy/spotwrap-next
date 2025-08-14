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
	"strings"
	"sync"

	"spotwrap-next/database"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// filenameFormat defines the output filename format
	filenameFormat = "{artist} - {title}.{output-ext}"
)

// Downloader handles downloading of tracks from Spotify
type Downloader struct {
	ctx        context.Context
	emitEvents bool
	db         *database.Database
}

// NewDownloader creates a new Downloader instance
func NewDownloader(db *database.Database) *Downloader {
	return &Downloader{db: db}
}

// Startup is called when the application starts
func (d *Downloader) Startup(ctx context.Context, emitEvents bool) {
	d.ctx = ctx
	d.emitEvents = emitEvents
}

// Download downloads a track from the provided Spotify link
// Parameters:
// - link: Spotify link to download from
// - outputPath: directory where to save the downloaded files
// - format: output format (mp3, wav, etc.)
// - bitrate: quality of the output (128k, 320k, etc.)
// - songsToCheck: optional list of songs to check after download
// Returns: boolean indicating whether the download was successful
func (d *Downloader) Download(link, outputPath, format, bitrate string, songsToCheck []string) bool {
	// Extract the spotdl binary to a temporary location
	tmpDir, err := os.MkdirTemp("", "spotdl")
	if err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to create temp directory: %v", err))
		return false
	}
	defer os.RemoveAll(tmpDir)

	// Determine file paths based on OS
	var spotdlPath string
	var ffmpegPath string
	isWindows := runtime.GOOS == "windows"

	if isWindows {
		if success := extractWindowsBinaries(d, tmpDir); !success {
			return false
		}
		spotdlPath = filepath.Join(tmpDir, "spotdl.exe")
		ffmpegPath = filepath.Join(tmpDir, "ffmpeg.exe")
	} else {
		if success := extractLinuxBinaries(d, tmpDir); !success {
			return false
		}
		spotdlPath = filepath.Join(tmpDir, "spotdl")
	}

	creds, err := d.db.GetSpotifyCredentials()
	if err != nil {
		log.Printf("Error fetching credentials: %v", err)
		return false
	}

	// Prepare arguments
	args := []string{
		link,
		"--bitrate", bitrate,
		"--format", format,
		"--client-id", creds.ClientID,
		"--client-secret", creds.ClientSecret,
		"--output",
	}

	outputFilePath := filenameFormat
	if outputPath != "" {
		outputFilePath = filepath.Join(outputPath, filenameFormat)
	}
	args = append(args, outputFilePath)

	// Add ffmpeg path argument for Windows
	if isWindows && ffmpegPath != "" {
		args = append(args, "--ffmpeg", ffmpegPath)
	}

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
	d.emitUpdateEvent("Downloading")
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

	// Optional post-download verification: if a list of expected files was provided
	// via songsToDelete, verify that each exists on disk. For any missing file,
	// emit an update event that the frontend can react to.
	if len(songsToCheck) > 0 {
		// Indicate verification start
		d.emitUpdateEvent("VerifyingDownload")

		for _, expectedFileName := range songsToCheck {
			var candidatePath string
			if outputPath != "" {
				candidatePath = filepath.Join(outputPath, expectedFileName)
			} else {
				candidatePath = expectedFileName
			}

			if _, statErr := os.Stat(candidatePath); os.IsNotExist(statErr) {
				// Prefix with a recognizable tag for the frontend to parse
				d.emitUpdateEvent(fmt.Sprintf("missing_track:%s", expectedFileName))
			}
		}
		d.emitUpdateEvent("VerificationComplete")
	}

	d.emitUpdateEvent("Done")
	return true
}

// extractBinary extracts a binary from the embedded FS to the target path
func (d *Downloader) extractBinary(embeddedFS embed.FS, embeddedPath, targetPath string) error {
	binData, err := embeddedFS.ReadFile(embeddedPath)
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

	// List of strings to remove from the output
	spotdl_strings := []string{
		"WARNING:root:",
		"INFO:spotdl.download.downloader:",
		"INFO:spotdl.utils.search:",
		"ERROR:spotipy.client:",
		"ERROR:spotdl.download.progress_handler:",
	}

	buf := make([]byte, 1024)
	for {
		n, err := pipe.Read(buf)
		if n > 0 {
			output := string(buf[:n])
			for _, spotdl_string := range spotdl_strings {
				output = strings.ReplaceAll(output, spotdl_string, "")
			}
			log.Print(output)
			d.emitUpdateEvent(output)
		}
		if err != nil {
			break
		}
	}
}

// emitErrorEvent emits a fatal error event and a done event
func (d *Downloader) emitErrorEvent(errMsg string) {
	d.emitUpdateEvent(fmt.Sprintf("fatal_error: %s", errMsg))
	d.emitUpdateEvent("Done")
}

// emitUpdateEvent emits an update event
func (d *Downloader) emitUpdateEvent(errMsg string) {
	if !d.emitEvents {
		return
	}
	wailsruntime.EventsEmit(d.ctx, "update_in_download", errMsg)
}
