package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"spotwrap-next/autostart"
	"spotwrap-next/spotdl"
	"spotwrap-next/utils"
	"syscall"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	noGUI := flag.Bool("no-gui", false, "Run in background mode")
	flag.Parse()

	if *noGUI {
		runInBackground()
		return
	}

	if err := startGUI(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func startGUI() error {
	// Create an instance of the app structure
	app, err := NewApp()
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	utils := utils.New()
	downloader := spotdl.NewDownloader()
	autostartSvc := autostart.New("spotwrap-next", "Spotwrap Next")

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
			autostartSvc,
		},
		CSSDragProperty:          "widows",
		CSSDragValue:             "1",
		EnableDefaultContextMenu: false,
		OnShutdown: func(ctx context.Context) {
			app.Close()
			utils.CleanUp() // clean the cover directory
		},
	})

	if err != nil {
		return fmt.Errorf("error running application: %w", err)
	}

	return nil
}

func runInBackground() {
	log.Println("Running in background mode")

	app, err := NewApp()
	if err != nil {
		log.Fatalf("Error initializing app: %v", err)
	}

	app.startup(context.Background())
	app.checkForNewReleases()
	app.startBackgroundChecker()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh // Wait for termination signal

	log.Println("Shutting down background service")
	app.stopBackgroundChecker()
}
