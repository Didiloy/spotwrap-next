package main

import (
	"context"
	"fmt"
	"spotwrap-next/api"
	"spotwrap-next/database"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                 context.Context
	spotifyAccessToken  string
	tokenExpirationTime time.Time
	database            *database.Database
}

// NewApp creates a new App application struct
func NewApp() (*App, error) {
	app := &App{}
	database, errDB := database.New()
	app.database = database
	if errDB != nil {
		fmt.Printf("Could not initialize database: \n%s\n", errDB.Error())
		return nil, errDB
	}
	go app.refreshTokenPeriodically() // Start automatic token refresh
	return app, nil
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.fetchSpotifyAccessToken()
}

// ================ Spotify =================
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

// ================ Database ==============
func (a *App) AddArtist(spotifyID string) bool {
	success, err := a.database.AddArtist(spotifyID)
	if err != nil {
		fmt.Println("Error adding artist:", err)
		return false
	}
	return success
}

func (a *App) RemoveArtist(spotifyID string) bool {
	success, err := a.database.RemoveArtist(spotifyID)
	if err != nil {
		fmt.Println("Error removing artist:", err)
		return false
	}
	return success
}

func (a *App) GetArtistsFromDB() []database.Artist {
	artists, err := a.database.GetArtistsFromDB()
	if err != nil {
		fmt.Println("Error getting artists:", err)
		return nil
	}
	return artists
}

// ================ Utils =================
func (a *App) ChooseDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Directory",
	})
	if err != nil {
		return ""
	}
	return dir
}

// func IsThereNewRelease(releases []map[string]any, date time.Time) (bool, map[string]any) {

// 	items, ok := albums["items"].([]any)
// 	if !ok {
// 		return false, nil
// 	}

// 	// Find the most recent release after the given date
// 	var newestRelease map[string]any
// 	hasNewRelease := false

// 	for _, item := range items {
// 		album, ok := item.(map[string]any)
// 		if !ok {
// 			continue
// 		}

// 		releaseDateStr, ok := album["release_date"].(string)
// 		if !ok {
// 			continue
// 		}

// 		// Parse the release date (format can be "YYYY", "YYYY-MM", or "YYYY-MM-DD")
// 		var releaseDate time.Time
// 		switch len(releaseDateStr) {
// 		case 4: // YYYY
// 			releaseDate, _ = time.Parse("2006", releaseDateStr)
// 		case 7: // YYYY-MM
// 			releaseDate, _ = time.Parse("2006-01", releaseDateStr)
// 		default: // YYYY-MM-DD
// 			releaseDate, _ = time.Parse("2006-01-02", releaseDateStr)
// 		}

// 		if releaseDate.After(date) {
// 			if !hasNewRelease || releaseDate.After(newestRelease["release_date"].(time.Time)) {
// 				newestRelease = album
// 				newestRelease["release_date"] = releaseDate
// 				hasNewRelease = true
// 			}
// 		}
// 	}

// 	return hasNewRelease, newestRelease
// }

func (a *App) Close() {
	a.database.Close()
}
