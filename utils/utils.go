package utils

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/EdlinOrg/prominentcolor"
)

type Utils struct {
	counter int
}

func New() *Utils {
	return &Utils{}
}

func (u *Utils) getCounter() int {
	u.counter++
	return u.counter
}

func (u *Utils) GetDominantColor(imageLink string) []string {
	colors, err := u.getDominantColor(imageLink)
	if err != nil {
		fmt.Printf("Could not get dominant colors for image %v: %v\n", imageLink, err)
		return make([]string, 0)
	}
	return colors
}

// GetDominantColor extracts the most dominant colors from an image URL
func (u *Utils) getDominantColor(imageLink string) ([]string, error) {
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

	os.MkdirAll("album_cover", os.ModePerm)
	filename := "album_cover/cover" + strconv.Itoa(u.getCounter())
	//forced to write the image file, otherwise it will not work
	err = os.WriteFile(filename, imgData, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to write image file: %w", err)
	}

	img, err := loadImage(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to load image: %w", err)
	}

	colours, err := prominentcolor.Kmeans(img)
	if err != nil {
		os.Remove(filename)
		return nil, fmt.Errorf("Failed to process image: %v", err)
	}

	var colors []string
	for _, colour := range colours {
		colors = append(colors, "#"+colour.AsString())
	}

	//delete the image file
	err = os.Remove(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to delete image file: %w", err)
	}

	return colors, nil
}

func (u *Utils) CleanUp() {
	fmt.Println("Deleting album cover directory...")
	os.RemoveAll("album_cover")
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
