package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"spotwrap-next/autostart"
	"spotwrap-next/spotdl"
	"spotwrap-next/utils"

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

	cliMode := flag.Bool("cli", false, "Run in background mode")
	flag.Parse()

	if *cliMode {
		runInBackground()
		os.Exit(0)
	}

	// Normal GUI startup
	err := startGui()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func startGui() error {
	// Create an instance of the app structure
	app, err := NewApp()
	if err != nil {
		return fmt.Errorf("Could not initialize app: \n%s\n", err.Error())
	}

	utils := utils.New()
	downloader := spotdl.NewDownloader()
	autostart := autostart.New("spotwrap-next", "Spotwrap Next")

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "spotwrap-next",
		Width:  1100,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			downloader.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			utils,
			downloader,
			autostart,
		},
		CSSDragProperty:          "widows",
		CSSDragValue:             "1",
		EnableDefaultContextMenu: false,
		OnShutdown: func(ctx context.Context) {
			app.Close()
			utils.CleanUp() //clean the cover directory
		},
	})

	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}

	return nil
}

func runInBackground() {
	fmt.Println("Running in background mode")
}
