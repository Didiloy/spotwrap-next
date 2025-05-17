//go:build windows
// +build windows

package spotdl

import (
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed assets/windows/spotdl.exe assets/windows/ffmpeg.exe assets/windows/ffplay.exe assets/windows/ffprobe.exe
var windowsBinaries embed.FS

// windowsExtractFunc extracts Windows-specific binaries to the target directory
func windowsExtractFunc(d *Downloader, tmpDir string) bool {
	// Define paths for Windows binaries
	spotdlPath := filepath.Join(tmpDir, "spotdl.exe")
	ffmpegPath := filepath.Join(tmpDir, "ffmpeg.exe")
	ffplayPath := filepath.Join(tmpDir, "ffplay.exe")
	ffprobePath := filepath.Join(tmpDir, "ffprobe.exe")

	// Extract spotdl.exe
	if err := d.extractBinary(windowsBinaries, "assets/windows/spotdl.exe", spotdlPath); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to extract spotdl binary: %v", err))
		return false
	}

	// Extract FFmpeg binaries
	if err := d.extractBinary(windowsBinaries, "assets/windows/ffmpeg.exe", ffmpegPath); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to extract ffmpeg binary: %v", err))
		return false
	}
	if err := d.extractBinary(windowsBinaries, "assets/windows/ffplay.exe", ffplayPath); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to extract ffplay binary: %v", err))
		return false
	}
	if err := d.extractBinary(windowsBinaries, "assets/windows/ffprobe.exe", ffprobePath); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to extract ffprobe binary: %v", err))
		return false
	}

	return true
}

func init() {
	extractWindowsBinaries = windowsExtractFunc
}
