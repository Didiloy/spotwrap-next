package utils

import (
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/EdlinOrg/prominentcolor"
)

// GetDominantColor extracts the most dominant colors from an image URL
func GetDominantColor(imageLink string) ([]string, error) {
	// Download the image with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
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

	filename := "album_cover"
	//forced to write the image file, otherwise it will not work
	os.WriteFile(filename, imgData, 0666)

	img, err := loadImage(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to load image: %w", err)
	}

	colours, err := prominentcolor.Kmeans(img)
	if err != nil {
		log.Fatal("Failed to process image", err)
	}

	var colors []string
	for _, colour := range colours {
		colors = append(colors, "#"+colour.AsString())
	}

	//delete the image file
	os.Remove(filename)

	return colors, nil
}

func loadImage(fileInput string) (image.Image, error) {
	f, err := os.Open(fileInput)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
