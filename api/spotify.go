package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	TokenURL  = "https://accounts.spotify.com/api/token"
	BaseURL   = "https://api.spotify.com/v1"
	SearchURL = BaseURL + "/search"
	ArtistURL = BaseURL + "/artists"
	AlbumURL  = BaseURL + "/albums"
	TrackURL  = BaseURL + "/tracks"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetToken(clientID, clientSecret string) (string, int, error) {
	if clientID == "" || clientSecret == "" {
		return "", 0, fmt.Errorf("Missing Spotify client ID or client secret")
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

func GetArtistDetails(id string, token string) (map[string]any, error) {
	artistData := make(map[string]any)

	// Get basic artist info using ArtistURL
	basicInfoURL := fmt.Sprintf("%s/%s", ArtistURL, id)
	basicInfo, err := makeRequest(basicInfoURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic artist info: %v", err)
	}
	artistData["artist"] = basicInfo

	// Get artist's top tracks
	topTracksURL := fmt.Sprintf("%s/%s/top-tracks?market=US", ArtistURL, id)
	topTracks, err := makeRequest(topTracksURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get top tracks: %v", err)
	}
	artistData["top_tracks"] = topTracks["tracks"]

	// Get artist's albums
	albumsURL := fmt.Sprintf("%s/%s/albums?include_groups=album,single&market=US&limit=10", ArtistURL, id)
	albums, err := makeRequest(albumsURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get albums: %v", err)
	}
	artistData["albums"] = albums["items"]

	return artistData, nil
}

func GetAlbumDetails(id string, token string) (map[string]any, error) {
	albumData := make(map[string]any)

	// Get basic album info
	albumURL := fmt.Sprintf("%s/albums/%s", BaseURL, id)
	albumInfo, err := makeRequest(albumURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get album info: %v", err)
	}
	albumData["album"] = albumInfo

	// Get album tracks
	tracksURL := fmt.Sprintf("%s/albums/%s/tracks?limit=50", BaseURL, id)
	tracks, err := makeRequest(tracksURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get album tracks: %v", err)
	}
	albumData["tracks"] = tracks["items"]

	return albumData, nil
}

func GetTrackDetails(id string, token string) (map[string]any, error) {
	trackData := make(map[string]any)

	// Get basic track info
	trackURL := fmt.Sprintf("%s/%s", TrackURL, id)
	trackInfo, err := makeRequest(trackURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get track info: %v", err)
	}
	trackData["track"] = trackInfo

	return trackData, nil
}

func makeRequest(url string, token string) (map[string]any, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
