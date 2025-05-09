// Package utils provides utility functions for the application
package utils

import (
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/EdlinOrg/prominentcolor"
)

const (
	coverDirName = "album_cover"
	timeout      = 10 * time.Second
)

// Utils provides utility functions for the application
type Utils struct {
	counter int
}

// New creates a new Utils instance
func New() *Utils {
	return &Utils{}
}

// incrementCounter increments and returns the counter value
func (u *Utils) incrementCounter() int {
	u.counter++
	return u.counter
}

// GetDominantColor extracts the dominant colors from an image URL
// It returns a slice of hex color codes representing the most prominent colors
func (u *Utils) GetDominantColor(imageLink string) []string {
	colors, err := u.extractDominantColors(imageLink)
	if err != nil {
		log.Printf("Could not get dominant colors for image %v: %v\n", imageLink, err)
		return []string{}
	}
	return colors
}

// extractDominantColors downloads an image and extracts its dominant colors
func (u *Utils) extractDominantColors(imageLink string) ([]string, error) {
	// Download the image with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(imageLink)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	// Read the image data into memory
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Ensure the cover directory exists
	if err := os.MkdirAll(coverDirName, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create cover directory: %w", err)
	}

	// Create a temporary file for the image
	filename := filepath.Join(coverDirName, fmt.Sprintf("cover%d", u.incrementCounter()))
	if err := os.WriteFile(filename, imgData, 0666); err != nil {
		return nil, fmt.Errorf("failed to write image file: %w", err)
	}

	// Ensure the file is deleted when we're done
	defer os.Remove(filename)

	img, err := loadImage(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}

	colours, err := prominentcolor.Kmeans(img)
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w", err)
	}

	var colors []string
	for _, colour := range colours {
		colors = append(colors, "#"+colour.AsString())
	}

	return colors, nil
}

// CleanUp removes the album cover directory and its contents
func (u *Utils) CleanUp() {
	log.Println("Cleaning up album cover directory...")
	if err := os.RemoveAll(coverDirName); err != nil {
		log.Printf("Failed to remove album cover directory: %v", err)
	}
}

// loadImage loads an image from a file path
func loadImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}
