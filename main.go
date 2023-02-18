package main

import (
	"bookmarks_v3/module/cfg"
	"bookmarks_v3/module/db"
	"bookmarks_v3/module/tui"
	"fmt"
	"os"
)

func main() {
	bookmarksPrefs := cfg.LoadConfig()

	_, err := os.Stat(bookmarksPrefs.Database)
	if err != nil {
		fmt.Println("Creating database")
		db.CreateDatabase(bookmarksPrefs.Database)
	}

	tui.MainLoop(bookmarksPrefs)
}
