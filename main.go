package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	appinternal "stax-desktop/internal/app"
	"stax-desktop/internal/database"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dbPath := resolveDBPath()
	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("échec de l'ouverture de la base de données: %v", err)
	}
	defer db.Close()

	app := appinternal.NewApp(db)

	err = wails.Run(&options.App{
		Title:  "STAX — Gestion de Stock & Point de Vente",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:        app.Startup,
		Bind:             []interface{}{app},
	})
	if err != nil {
		log.Fatalf("erreur Wails: %v", err)
	}
}

// resolveDBPath puts the database in the user's AppData folder on Windows
// (or $HOME/.stax on other OSes during development), so it survives app updates.
func resolveDBPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	dir := filepath.Join(configDir, "STAX")
	_ = os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "stax.db")
}
