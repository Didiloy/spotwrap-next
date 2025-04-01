package main

import (
	"context"
	"fmt"
	"spotwrap-next/api"
	"spotwrap-next/utils"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                 context.Context
	spotifyAccessToken  string
	tokenExpirationTime time.Time
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	go app.refreshTokenPeriodically() // Start automatic token refresh
	return app
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.fetchSpotifyAccessToken()
}

// Fetch Spotify Access Token if expired
func (a *App) fetchSpotifyAccessToken() {
	if time.Now().After(a.tokenExpirationTime) { // Check if token is expired
		token, expiresIn, err := api.GetToken()
		if err != nil {
			fmt.Println("Error fetching token:", err)
			return
		}
		a.spotifyAccessToken = token
		a.tokenExpirationTime = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}
}

// Goroutine to refresh token every 55 minutes
func (a *App) refreshTokenPeriodically() {
	ticker := time.NewTicker(55 * time.Minute) // Refresh 5 min before expiry
	defer ticker.Stop()

	for {
		<-ticker.C
		a.fetchSpotifyAccessToken()
		fmt.Println("Token refreshed")
	}
}

// Search Spotify API for query
func (a *App) Search(query string) map[string]any {
	result, err := api.Search(query, a.spotifyAccessToken)
	if err != nil {
		fmt.Println("Error searching:", err)
		return map[string]any{}
	}
	return result
}

// Get Artist Data
func (a *App) GetArtist(id string) map[string]any {
	result, err := api.GetArtistDetails(id, a.spotifyAccessToken)
	if err != nil {
		fmt.Println("Error getting artist:", err)
		return map[string]any{}
	}
	return result
}

// Get Album Data
func (a *App) GetAlbum(id string) map[string]any {
	result, err := api.GetAlbumDetails(id, a.spotifyAccessToken)
	if err != nil {
		fmt.Println("Error getting album:", err)
		return map[string]any{}
	}
	return result
}

// Get Track Data
func (a *App) GetTrack(id string) map[string]any {
	result, err := api.GetTrackDetails(id, a.spotifyAccessToken)
	if err != nil {
		fmt.Println("Error getting Track:", err)
		return map[string]any{}
	}
	return result
}

func (a *App) GetDominantColor(imageLink string) ([]string, error) {
	colors, err := utils.GetDominantColor(imageLink)
	if err != nil {
		fmt.Printf("Could not get dominant colors for image %v: %v\n", imageLink, err)
		return []string{}, nil
	}
	return colors, nil
}

func (a *App) ChooseDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Directory",
	})
	if err != nil {
		return ""
	}
	return dir
}
