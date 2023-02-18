package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	ui "github.com/jairochavesb/clui"
	_ "github.com/mattn/go-sqlite3"
)

func SearchByID(id, database string) string {
	searchQuery := "SELECT * FROM bookmarks where id=" + id + ";"

	sqlDB, err := sql.Open("sqlite3", database)
	if err != nil {
		return err.Error()
	}

	defer sqlDB.Close()

	row, err := sqlDB.Query(searchQuery)
	if err != nil {
		return err.Error()
	}
	defer row.Close()

	var id2 string
	var name string
	var url string
	var tags string

	for row.Next() {
		row.Scan(&id2, &name, &tags, &url)
	}

	return name + "■" + tags + "■" + url
}

func SearchBookmark(listbox *ui.ListBox, field, keywords, database string) string {
	listbox.Clear()

	keywordsSeparated := strings.SplitAfter(keywords, " ")

	searchQuery := "SELECT * FROM bookmarks where " + field + " "

	for _, keyword := range keywordsSeparated {

		keyword = strings.ReplaceAll(keyword, " ", "")

		if keyword == "AND" {
			searchQuery += " AND " + field + ""

		} else if keyword == "OR" {
			searchQuery += " OR " + field + " "

		} else {
			searchQuery += "like \"%" + keyword + "%\" "
		}
	}

	searchQuery += ";"

	sqlDB, err := sql.Open("sqlite3", database)
	if err != nil {
		return err.Error()
	}

	defer sqlDB.Close()
	row, err := sqlDB.Query(searchQuery)
	if err != nil {
		return err.Error()
	}

	defer row.Close()

	var id string
	var name string
	var url string
	var tags string

	counter := 0

	for row.Next() {
		row.Scan(&id, &name, &tags, &url)

		if len(name) > 45 {
			tmp := name[0:45]
			name = tmp + "..."
		}

		if len(tags) > 55 {
			tmp := tags[0:55]
			tags = tmp + "..."
		}

		if len(url) > 55 {
			tmp := url[0:55]
			url = tmp + "..."
		}

		idNameDiv := strings.Repeat(" ", 12-len(id))
		nameTagsDiv := strings.Repeat(" ", 49-len(name))
		tagsUrlDiv := strings.Repeat(" ", 59-len(tags))

		listbox.AddItem(id + idNameDiv + name + nameTagsDiv + tags + tagsUrlDiv + url)

		counter++
	}

	return fmt.Sprintf("Search returned %d items", counter)
}

func DeleteBookmark(id, database string) string {
	deleteQuery := "DELETE FROM bookmarks WHERE id=" + id + ";"

	sqlDB, err := sql.Open("sqlite3", database)
	if err != nil {
		return err.Error()
	}
	defer sqlDB.Close()

	statement, err := sqlDB.Prepare(deleteQuery)
	if err != nil {
		return err.Error()
	}

	_, err = statement.Exec()
	if err != nil {
		return err.Error()
	}

	return "Bookmark updated successfully"
}

func UpdateBookmark(id, name, tags, url, database string) string {

	updateQuery := "UPDATE bookmarks set name=\"" + name + "\","
	updateQuery += "tags=\"" + tags + "\","
	updateQuery += "url=\"" + url + "\" where id=\"" + id + "\";"

	sqlDB, err := sql.Open("sqlite3", database)
	if err != nil {
		return err.Error()
	}
	defer sqlDB.Close()

	statement, err := sqlDB.Prepare(updateQuery)
	if err != nil {
		return err.Error()
	}

	_, err = statement.Exec()
	if err != nil {
		return err.Error()
	}

	return "Bookmark updated successfully"
}

func InsertIntoDB(name string, tags string, url string, database string) string {
	q := "SELECT COUNT(url) FROM bookmarks WHERE url=\"" + url + "\";"

	sqlDB, err := sql.Open("sqlite3", database)
	if err != nil {
		return err.Error()
	}
	defer sqlDB.Close()

	row, err := sqlDB.Query(q)
	if err != nil {
		return err.Error()
	}

	defer row.Close()
	var count int

	for row.Next() {
		row.Scan(&count)
	}

	if count >= 1 {
		return "Bookmark already exist"
	}

	insert := "INSERT INTO bookmarks(name, tags, url) VALUES ("
	insert += "'" + name + "',"
	insert += "'" + tags + "',"
	insert += "'" + url + "');"

	sqlDB, err = sql.Open("sqlite3", database)
	if err != nil {
		return err.Error()
	}
	defer sqlDB.Close()

	statement, err := sqlDB.Prepare(insert)
	if err != nil {
		return err.Error()
	}

	_, err = statement.Exec()
	if err != nil {
		return err.Error()
	}

	return "Bookmark added successfully!"
}

func CreateDatabase(database string) {
	file, err := os.Create(database)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	sqlDB, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}

	createItemsTable := `CREATE TABLE bookmarks (
      "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,      
      "name" TEXT,
      "tags" TEXT,
      "url" TEXT
     );`

	statement, err := sqlDB.Prepare(createItemsTable)
	if err != nil {
		log.Fatal(err)
	}

	statement.Exec()
	defer sqlDB.Close()
}
