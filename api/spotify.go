package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const (
	TokenURL  = "https://accounts.spotify.com/api/token"
	BaseURL   = "https://api.spotify.com/v1"
	SearchURL = BaseURL + "/search"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetToken() (string, int, error) {
	err := godotenv.Load()
	if err != nil {
		return "", 0, fmt.Errorf("Error loading .env file: %v", err)
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", 0, fmt.Errorf("Missing SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET")
	}

	// Prepare request body
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	// Decode response
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", 0, err
	}

	// Log token expiration
	// fmt.Println("Spotify Token Expires In (seconds):", strconv.Itoa(tokenResponse.ExpiresIn))

	return tokenResponse.AccessToken, tokenResponse.ExpiresIn, nil
}

func Search(query string, token string) (map[string]any, error) {
	req, err := http.NewRequest("GET", SearchURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("type", "album,artist,track")
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
