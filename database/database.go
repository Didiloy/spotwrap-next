package database

import (
	"database/sql"
	"fmt"
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

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS artists (
		spotify_id TEXT PRIMARY KEY,
		last_checked TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Printf("Database initialized at %s\n", dbPath)
	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	fmt.Println("Closing database connection...")
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
	rows, err := d.db.Query("SELECT spotify_id, last_checked FROM artists")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []Artist
	for rows.Next() {
		var a Artist
		if err := rows.Scan(&a.SpotifyID, &a.LastChecked); err != nil {
			return nil, err
		}
		artists = append(artists, a)
	}
	return artists, nil
}
