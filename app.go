package main

import (
	"context"
	"fmt"
	"log"
	"spotwrap-next/api"
	"spotwrap-next/database"
	"spotwrap-next/notifications"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App represents the main application structure
type App struct {
	ctx                 context.Context
	spotifyAccessToken  string
	tokenExpirationTime time.Time
	db                  *database.Database
	backgroundTicker    *time.Ticker
	backgroundDone      chan bool
}

// NewApp creates a new App application struct
func NewApp() (*App, error) {
	db, err := database.New()
	if err != nil {
		return nil, fmt.Errorf("database initialization failed: %w", err)
	}

	app := &App{
		db:             db,
		backgroundDone: make(chan bool),
	}

	// Start automatic token refresh in a goroutine
	go app.refreshTokenPeriodically()

	return app, nil
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.fetchSpotifyAccessToken()
}

// fetchSpotifyAccessToken retrieves and updates the Spotify access token if it has expired
func (a *App) fetchSpotifyAccessToken() {
	if time.Now().Before(a.tokenExpirationTime) {
		// Token still valid
		return
	}

	// Get credentials from database
	creds, err := a.db.GetSpotifyCredentials()
	if err != nil {
		log.Printf("Error fetching credentials: %v", err)
		return
	}

	token, expiresIn, err := api.GetToken(creds.ClientID, creds.ClientSecret)
	if err != nil {
		log.Printf("Error fetching token: %v", err)
		return
	}

	a.spotifyAccessToken = token
	a.tokenExpirationTime = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// refreshTokenPeriodically refreshes the Spotify token every 55 minutes
func (a *App) refreshTokenPeriodically() {
	ticker := time.NewTicker(55 * time.Minute) // Refresh 5 min before expiry
	defer ticker.Stop()

	for range ticker.C {
		a.fetchSpotifyAccessToken()
		log.Println("Token refreshed")
	}
}

// Search queries the Spotify API with the given query string
func (a *App) Search(query string) map[string]any {
	result, err := api.Search(query, a.spotifyAccessToken)
	if err != nil {
		log.Printf("Error searching: %v", err)
		return map[string]any{}
	}
	return result
}

// GetArtist retrieves artist data from Spotify by ID
func (a *App) GetArtist(id string) map[string]any {
	result, err := api.GetArtistDetails(id, a.spotifyAccessToken, false)
	if err != nil {
		log.Printf("Error getting artist: %v", err)
		return map[string]any{}
	}
	return result
}

// GetAlbum retrieves album data from Spotify by ID
func (a *App) GetAlbum(id string) map[string]any {
	result, err := api.GetAlbumDetails(id, a.spotifyAccessToken)
	if err != nil {
		log.Printf("Error getting album: %v", err)
		return map[string]any{}
	}
	return result
}

// GetTrack retrieves track data from Spotify by ID
func (a *App) GetTrack(id string) map[string]any {
	result, err := api.GetTrackDetails(id, a.spotifyAccessToken)
	if err != nil {
		log.Printf("Error getting track: %v", err)
		return map[string]any{}
	}
	return result
}

// AddArtist adds an artist to the database by Spotify ID
func (a *App) AddArtist(spotifyID string) bool {
	success, err := a.db.AddArtist(spotifyID)
	if err != nil {
		log.Printf("Error adding artist: %v", err)
		return false
	}
	return success
}

// RemoveArtist removes an artist from the database by Spotify ID
func (a *App) RemoveArtist(spotifyID string) bool {
	success, err := a.db.RemoveArtist(spotifyID)
	if err != nil {
		log.Printf("Error removing artist: %v", err)
		return false
	}
	return success
}

// GetArtistsFromDB retrieves all artists from the database
func (a *App) GetArtistsFromDB() []database.Artist {
	artists, err := a.db.GetArtistsFromDB()
	if err != nil {
		log.Printf("Error getting artists: %v", err)
		return nil
	}
	return artists
}

// ================ Spotify Credentials =================
// GetSpotifyCredentials retrieves Spotify credentials from the database
func (a *App) GetSpotifyCredentials() map[string]string {
	creds, err := a.db.GetSpotifyCredentials()
	if err != nil {
		log.Printf("Error getting Spotify credentials: %v", err)
		return map[string]string{
			"clientId":     "",
			"clientSecret": "",
		}
	}

	return map[string]string{
		"clientId":     creds.ClientID,
		"clientSecret": creds.ClientSecret,
	}
}

// SetSpotifyCredentials validates and stores Spotify credentials
func (a *App) SetSpotifyCredentials(clientID, clientSecret string) bool {
	// First check if the credentials are valid by trying to get a token
	token, _, err := api.GetToken(clientID, clientSecret)
	if err != nil || token == "" {
		return false
	}

	// Store credentials in database
	if err := a.db.StoreSpotifyCredentials(clientID, clientSecret); err != nil {
		log.Printf("Error storing Spotify credentials: %v", err)
		return false
	}

	// If credentials are valid, refresh the token immediately
	a.fetchSpotifyAccessToken()
	return true
}

// HasValidSpotifyCredentials checks if the stored credentials are valid
func (a *App) HasValidSpotifyCredentials() bool {
	creds, err := a.db.GetSpotifyCredentials()
	if err != nil {
		log.Printf("Error getting Spotify credentials: %v", err)
		return false
	}

	token, _, err := api.GetToken(creds.ClientID, creds.ClientSecret)
	return err == nil && token != ""
}

// ChooseDirectory opens a directory selection dialog
func (a *App) ChooseDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Directory",
	})
	if err != nil {
		return ""
	}
	return dir
}

// ================ Utils =================
func (a *App) IsANewRelease(id string, release map[string]any) bool {
	artist, err := a.db.GetArtistByID(id)
	if err != nil {
		fmt.Println("Error getting artist:", err)
		return false
	}

	// Extract release date from the album
	releaseDateStr, ok := release["release_date"].(string)
	if !ok {
		return false
	}

	// Parse the release date (handles different formats: YYYY, YYYY-MM, YYYY-MM-DD)
	var releaseDate time.Time
	switch len(releaseDateStr) {
	case 4: // YYYY
		releaseDate, _ = time.Parse("2006", releaseDateStr)
	case 7: // YYYY-MM
		releaseDate, _ = time.Parse("2006-01", releaseDateStr)
	default: // YYYY-MM-DD
		releaseDate, _ = time.Parse("2006-01-02", releaseDateStr)
	}

	return releaseDate.After(artist.LastChecked)
}

// Background
func (a *App) startBackgroundChecker() {
	a.backgroundTicker = time.NewTicker(5 * time.Hour)
	a.backgroundDone = make(chan bool)

	go func() {
		for {
			select {
			case <-a.backgroundTicker.C:
				a.checkForNewReleases()
			case <-a.backgroundDone:
				return
			}
		}
	}()
}

func (a *App) stopBackgroundChecker() {
	if a.backgroundTicker != nil {
		a.backgroundTicker.Stop()
	}
	if a.backgroundDone != nil {
		a.backgroundDone <- true
	}
}

func (a *App) checkForNewReleases() {
	fmt.Println("Starting background check for new releases...")

	// Get artists that need checking
	artists, err := a.db.GetArtistsFromDB()
	if err != nil {
		fmt.Printf("Error getting artists to check: %v\n", err)
		return
	}

	if len(artists) == 0 {
		fmt.Println("No artists need checking at this time")
		return
	}

	a.fetchSpotifyAccessToken()

	for _, artist := range artists {
		fmt.Printf("Checking for new releases from artist %s...\n", artist.SpotifyID)

		// Get artist's latest albums with retry mechanism
		artistData, err := api.GetArtistDetails(artist.SpotifyID, a.spotifyAccessToken, true)
		if err != nil {
			fmt.Printf("Error getting artist details for %s: %v\n", artist.SpotifyID, err)
			continue
		}

		// Check albums for new releases
		albums, ok := artistData["albums"].([]any)
		if !ok {
			fmt.Printf("Unexpected albums format for artist %s\n", artist.SpotifyID)
			continue
		}

		for _, album := range albums {
			albumMap, ok := album.(map[string]any)
			if !ok {
				continue
			}

			if a.IsANewRelease(artist.SpotifyID, albumMap) {
				fmt.Printf("New release found for artist %s: %v\n", artist.SpotifyID, albumMap["name"])

				albumName := albumMap["name"].(string)
				artistName := artistData["artist"].(map[string]any)["name"].(string)
				message := fmt.Sprintf("%s has released %s", artistName, albumName)

				// Send desktop notification
				err := notifications.Notify("New Release!", message)
				if err != nil {
					fmt.Printf("Failed to send notification: %v\n", err)
				}
			}
		}

		// Update last checked time
		if _, err := a.db.AddArtist(artist.SpotifyID); err != nil {
			fmt.Printf("Error updating last_checked for artist %s: %v\n", artist.SpotifyID, err)
		}
	}

	fmt.Println("Background check completed")
}

func (a *App) Close() {
	a.db.Close()
}
