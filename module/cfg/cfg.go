package cfg

import (
	"bookmarks_v3/module/misc"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func LoadConfig() *misc.BookmarksPreferences {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	configFile := userHomeDir + "/.config/bookmarks.cfg"

	_, err = os.Stat(configFile)
	if err != nil {
		SetConfig(configFile)
	}

	data, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	scanner := bufio.NewScanner(data)

	bookmarksPref := misc.BookmarksPreferences{}

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "browser=") {
			bookmarksPref.Browser = strings.Split(scanner.Text(), "=")[1]

		} else if strings.Contains(scanner.Text(), "database=") {
			bookmarksPref.Database = strings.Split(scanner.Text(), "=")[1]
		}
	}

	return &bookmarksPref
}

func SetConfig(configFile string) {
	fmt.Printf("Please asnwer the following questions.\n")

	webBrowser := getBrowserName()

	bookmarksDbDir := getDatabaseDir()

	bookmarksDbname := ""
	fmt.Println("How to name the bookmarks db file, leave this blank to use default name 'bookmarks.db'")
	fmt.Scanln(&bookmarksDbname)

	if bookmarksDbname == "" {
		bookmarksDbname = "bookmarks.db"
	}

	bookmarksdb := bookmarksDbDir + bookmarksDbname

	fd, err := os.Create(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	configText := "browser=" + webBrowser + "\ndatabase=" + bookmarksdb + "\n"

	_, err = fd.Write([]byte(configText))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getBrowserName() string {
	browser := ""

	for {
		fmt.Println("Favorite web browser to open bookmarks: ")
		fmt.Scanln(&browser)

		_, err := os.Stat(browser)
		if err != nil {
			fmt.Println("Web browser path do not exist. Try again.")
			continue
		} else {
			break
		}
	}

	return browser
}

func getDatabaseDir() string {
	databaseDir := ""

	for {
		fmt.Println("Where to place the bookmarks database: ")
		fmt.Scanln(&databaseDir)

		_, err := os.Stat(databaseDir)
		if err != nil {
			fmt.Println("Directory do not exist. Try again.")
			continue
		} else {
			break
		}
	}

	if databaseDir[len(databaseDir)-1] != '/' {
		databaseDir += "/"
	}

	return databaseDir
}
