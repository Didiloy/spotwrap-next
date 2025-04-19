package spotdl

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const filenameFormat = "{artist} - {title}.{output-ext}"

//go:embed assets/spotdl_linux
var spotdlBinary embed.FS

type Downloader struct {
	ctx context.Context
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Startup(ctx context.Context) {
	d.ctx = ctx
}

func (d *Downloader) Download(link, outputPath, format, bitrate string, songsToDelete []string) error {
	// Extract the spotdl binary to a temporary location
	tmpDir, err := os.MkdirTemp("", "spotdl")
	if err != nil {
		runtime.EventsEmit(d.ctx, "update_in_download", fmt.Sprintf("fatal_error: %v", err))
		runtime.EventsEmit(d.ctx, "update_in_download", "Done")
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	spotdlPath := filepath.Join(tmpDir, "spotdl")
	// if runtime.GOOS == "windows" {
	// 	spotdlPath += ".exe"
	// }

	// Get the embedded binary
	binData, err := spotdlBinary.ReadFile("assets/spotdl_linux")
	if err != nil {
		runtime.EventsEmit(d.ctx, "update_in_download", fmt.Sprintf("fatal_error: %v", err))
		runtime.EventsEmit(d.ctx, "update_in_download", "Done")
		return fmt.Errorf("failed to read embedded binary: %v", err)
	}

	// Write the binary to the temp location
	err = os.WriteFile(spotdlPath, binData, 0755)
	if err != nil {
		runtime.EventsEmit(d.ctx, "update_in_download", fmt.Sprintf("fatal_error: %v", err))
		runtime.EventsEmit(d.ctx, "update_in_download", "Done")
		return fmt.Errorf("failed to write binary to temp location: %v", err)
	}

	// Prepare arguments
	args := []string{
		link,
		"--bitrate",
		bitrate,
		"--format",
		format,
		"--print-errors",
		"--output",
	}

	if outputPath != "" {
		args = append(args, filepath.Join(outputPath, filenameFormat))
	} else {
		args = append(args, filenameFormat)
	}

	// Execute spotdl
	cmd := exec.Command(spotdlPath, args...)

	// Set up stdout and stderr pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		runtime.EventsEmit(d.ctx, "update_in_download", fmt.Sprintf("fatal_error: %v", err))
		runtime.EventsEmit(d.ctx, "update_in_download", "Done")
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		runtime.EventsEmit(d.ctx, "update_in_download", fmt.Sprintf("fatal_error: %v", err))
		runtime.EventsEmit(d.ctx, "update_in_download", "Done")
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	runtime.EventsEmit(d.ctx, "update_in_download", "Downloading")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Create a wait group to wait for the goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	//Stdout
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := stdoutPipe.Read(buf)
			if n > 0 {
				output := string(buf[:n])
				fmt.Printf("%s", output)
				runtime.EventsEmit(d.ctx, "update_in_download", output)
			}
			if err != nil {
				break
			}
		}
	}()

	//Stderr
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := stderrPipe.Read(buf)
			if n > 0 {
				output := string(buf[:n])
				fmt.Printf("%s", output)
				runtime.EventsEmit(d.ctx, "update_in_download", output)
			}
			if err != nil {
				break
			}
		}
	}()

	// Wait for the command to finish
	err = cmd.Wait()
	wg.Wait()

	if err != nil {
		runtime.EventsEmit(d.ctx, "update_in_download", fmt.Sprintf("fatal_error: %v", err))
		runtime.EventsEmit(d.ctx, "update_in_download", "Done")
		return fmt.Errorf("command execution failed: %v", err)
	}

	// Delete songs if requested
	// if len(songsToDelete) > 0 {
	// 	for _, song := range songsToDelete {
	// 		path := filepath.Join(outputPath, fmt.Sprintf("%s.%s", song, format))
	// 		err := os.Remove(path)
	// 		if err != nil {
	// 			msg := fmt.Sprintf("could not delete song %s", song)
	// 			fmt.Println(msg)
	// 			runtime.EventsEmit(d.ctx, "update_in_download", msg)
	// 			continue
	// 		}
	// 		msg := fmt.Sprintf("%s was deleted", path)
	// 		fmt.Println(msg)
	// 		runtime.EventsEmit(d.ctx, "update_in_download", msg)
	// 	}
	// }

	runtime.EventsEmit(d.ctx, "update_in_download", "Done")
	return nil
}
