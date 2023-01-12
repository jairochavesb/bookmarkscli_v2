package main

import (
	"bookmarks_v3/module/cfg"
	"bookmarks_v3/module/db"
	"bookmarks_v3/module/tui"
	"fmt"
	"log"
	"os"
)

func main() {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	configPath := userHomeDir + "/.config/bookmarks.cfg"
	bookmarksPrefs := cfg.LoadConfig(configPath)

	_, err = os.Stat(configPath)
	if err != nil {
		fmt.Printf("Stat: %s\n", err.Error())
		cfg.SetConfig(configPath)
		db.CreateDatabase(bookmarksPrefs.Database)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "--set-db" {
			cfg.SetConfig(configPath)
			db.CreateDatabase(bookmarksPrefs.Database)
		}
	}

	_, err = os.Stat(bookmarksPrefs.Database)
	if err != nil {
		fmt.Println("Creating database")
		db.CreateDatabase(bookmarksPrefs.Database)
	}

	tui.MainLoop(bookmarksPrefs)

}
