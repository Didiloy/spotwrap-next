//go:build linux
// +build linux

package spotdl

import (
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed assets/spotdl_linux
var linuxBinary embed.FS

// linuxExtractFunc extracts Linux-specific binary to the target directory
func linuxExtractFunc(d *Downloader, tmpDir string) bool {
	// Define path for Linux binary
	spotdlPath := filepath.Join(tmpDir, "spotdl")
	
	// Extract spotdl binary
	if err := d.extractBinary(linuxBinary, "assets/spotdl_linux", spotdlPath); err != nil {
		d.emitErrorEvent(fmt.Sprintf("failed to extract spotdl binary: %v", err))
		return false
	}
	
	return true
}

func init() {
	extractLinuxBinaries = linuxExtractFunc
}