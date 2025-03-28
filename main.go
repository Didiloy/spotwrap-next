package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/joho/godotenv"
)

//go:embed .env
var envFile embed.FS

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	//env file
	if data, err := envFile.ReadFile(".env"); err == nil {
		envMap, err := godotenv.Unmarshal(string(data))
		if err == nil {
			for k, v := range envMap {
				os.Setenv(k, v)
			}
		}
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "spotwrap-next",
		Width:  1100,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
