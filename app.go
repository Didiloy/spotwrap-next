package main

import (
	"context"
	"fmt"
	"spotwrap-next/api"
	"time"
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
