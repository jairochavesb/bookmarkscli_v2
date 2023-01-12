package cfg

import (
	"bookmarks_v3/module/misc"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func LoadConfig(cfgDir string) *misc.BookmarksPreferences {
	_, err := os.Stat(cfgDir)
	if err != nil {
		SetConfig(cfgDir)
	}

	data, err := os.Open(cfgDir)
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

func SetConfig(cfgDir string) {
	fmt.Println("Please asnwer the following questions.")

	fmt.Println("Favorite web browser to open bookmarks: ")
	webBrowser := ""
	fmt.Scanln(&webBrowser)

	fmt.Println("Where to place the bookmarks database: ")
	bookmarksDb := ""
	fmt.Scanln(&bookmarksDb)

	fd, err := os.Create(cfgDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	configText := "browser=" + webBrowser + "\ndatabase=" + bookmarksDb + "\n"

	_, err = fd.Write([]byte(configText))
	if err != nil {
		log.Fatal(err.Error())
	}

}
