package spotdl

// extractWindowsBinaries extracts Windows-specific binaries to the target directory
// It's implemented in windows.go for Windows builds
var extractWindowsBinaries func(d *Downloader, tmpDir string) bool

// extractLinuxBinaries extracts Linux-specific binary to the target directory
// It's implemented in linux.go for Linux builds
var extractLinuxBinaries func(d *Downloader, tmpDir string) bool

func init() {
	// Set up stub functions for platforms where the actual implementation doesn't exist
	if extractWindowsBinaries == nil {
		extractWindowsBinaries = func(d *Downloader, tmpDir string) bool {
			// This should never be called on non-Windows platforms
			d.emitErrorEvent("Windows binary extraction attempted on non-Windows platform")
			return false
		}
	}
	
	if extractLinuxBinaries == nil {
		extractLinuxBinaries = func(d *Downloader, tmpDir string) bool {
			// This should never be called on non-Linux platforms
			d.emitErrorEvent("Linux binary extraction attempted on non-Linux platform")
			return false
		}
	}
}