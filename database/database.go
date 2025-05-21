package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

type Artist struct {
	SpotifyID   string
	LastChecked time.Time
	CreatedAt   time.Time
}

type SpotifyCredentials struct {
	ClientID     string
	ClientSecret string
}

// New creates and initializes a new Database instance
func New() (*Database, error) {
	// Get user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// Create app directory if it doesn't exist
	appDir := filepath.Join(configDir, "spotwrap-next")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app directory: %w", err)
	}

	dbPath := filepath.Join(appDir, "artists.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create artists table
	createArtistsTableSQL := `
	CREATE TABLE IF NOT EXISTS artists (
		spotify_id TEXT PRIMARY KEY,
		last_checked TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createArtistsTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create artists table: %w", err)
	}

	// Create settings table for Spotify credentials
	createSettingsTableSQL := `
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`

	if _, err := db.Exec(createSettingsTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create settings table: %w", err)
	}

	fmt.Printf("Database initialized at %s\n", dbPath)
	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	log.Println("Closing database connection...")
	return d.db.Close()
}

// AddArtist adds or updates an artist in the database
func (d *Database) AddArtist(spotifyID string) (bool, error) {
	_, err := d.db.Exec(`
		INSERT INTO artists (spotify_id, last_checked)
		VALUES (?, ?)
		ON CONFLICT(spotify_id)
		DO UPDATE SET last_checked = ?`,
		spotifyID, time.Now(), time.Now())
	if err != nil {
		return false, err
	}
	return true, nil
}

// RemoveArtist removes an artist from the database
func (d *Database) RemoveArtist(spotifyID string) (bool, error) {
	_, err := d.db.Exec("DELETE FROM artists WHERE spotify_id = ?", spotifyID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetArtists retrieves all subscribed artists
func (d *Database) GetArtistsFromDB() ([]Artist, error) {
	rows, err := d.db.Query("SELECT spotify_id, last_checked, created_at FROM artists")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []Artist
	for rows.Next() {
		var a Artist
		if err := rows.Scan(&a.SpotifyID, &a.LastChecked, &a.CreatedAt); err != nil {
			return nil, err
		}
		artists = append(artists, a)
	}

	if len(artists) == 0 {
		artists = make([]Artist, 0)
	}
	return artists, nil
}

func (d *Database) GetArtistByID(id string) (*Artist, error) {
	row := d.db.QueryRow("SELECT spotify_id, last_checked, created_at FROM artists WHERE spotify_id = ?", id)
	var a Artist
	if err := row.Scan(&a.SpotifyID, &a.LastChecked, &a.CreatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

// StoreSpotifyCredentials saves Spotify API credentials to the database
func (d *Database) StoreSpotifyCredentials(clientID, clientSecret string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	// Store client ID
	_, err = tx.Exec(
		"INSERT INTO settings (key, value) VALUES ('spotify_client_id', ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		clientID, clientID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Store client secret
	_, err = tx.Exec(
		"INSERT INTO settings (key, value) VALUES ('spotify_client_secret', ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		clientSecret, clientSecret,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetSpotifyCredentials retrieves Spotify API credentials from the database
func (d *Database) GetSpotifyCredentials() (SpotifyCredentials, error) {
	var creds SpotifyCredentials

	// Get client ID
	var clientID sql.NullString
	err := d.db.QueryRow("SELECT value FROM settings WHERE key = 'spotify_client_id'").Scan(&clientID)
	if err != nil && err != sql.ErrNoRows {
		return creds, err
	}
	if clientID.Valid {
		creds.ClientID = clientID.String
	}

	// Get client secret
	var clientSecret sql.NullString
	err = d.db.QueryRow("SELECT value FROM settings WHERE key = 'spotify_client_secret'").Scan(&clientSecret)
	if err != nil && err != sql.ErrNoRows {
		return creds, err
	}
	if clientSecret.Valid {
		creds.ClientSecret = clientSecret.String
	}

	return creds, nil
}

// SetSetting saves a key-value pair to the settings table.
func (d *Database) SetSetting(key string, value string) error {
	_, err := d.db.Exec(
		"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		key, value, value,
	)
	return err
}

// GetSetting retrieves a value from the settings table by its key.
// It returns an empty string and no error if the key is not found.
func (d *Database) GetSetting(key string) (string, error) {
	var val sql.NullString
	err := d.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&val)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Key not found, return empty string and no error
		}
		return "", err // Other error
	}
	if val.Valid {
		return val.String, nil
	}
	return "", nil // Should not happen if ErrNoRows is handled, but as a fallback
}
