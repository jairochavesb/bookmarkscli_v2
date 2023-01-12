package tui

import (
	"bookmarks_v3/module/db"
	"bookmarks_v3/module/misc"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	ui "github.com/jairochavesb/clui"
	termbox "github.com/nsf/termbox-go"

	"golang.org/x/term"
)

type WindowSize struct {
	screenWidth  int
	screeHeight  int
	windowWidth  int
	windowHeight int
	windowX      int
	windowY      int
}

func MainLoop(bookmarksPrefs *misc.BookmarksPreferences) {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	_, err := os.Stat("themes/simpleDark.theme")
	if err == nil {
		ui.SetThemePath("themes/")
		ui.SetCurrentTheme("simpleDark")
	}

	initUI(bookmarksPrefs)

	ui.MainLoop()
}

func initUI(bookmarksPrefs *misc.BookmarksPreferences) {

	// Some calculations to make the option windows to appear in the center of the screen
	winSize := WindowSize{}

	winSize.screenWidth, winSize.screeHeight, _ = term.GetSize(0)

	winSize.windowWidth = 80
	winSize.windowHeight = 10

	winSize.windowY = (winSize.screeHeight / 2) - (winSize.windowHeight / 2)
	winSize.windowX = (winSize.screenWidth / 2) - (winSize.windowWidth / 2)

	// Start to draw the main window
	mainWindow := ui.AddWindow(0, 0, winSize.screenWidth, winSize.screeHeight, "")
	mainWindow.SetBorder(2)
	mainWindow.SetTitleButtons(0)
	mainWindow.SetTitle(" BOOKMARKS ")
	mainWindow.SetPack(ui.Vertical)

	menu1 := "<t: cyan bold>F1: <t: white bold>New\t\t\t<f: cyan bold>F2:<t: white bold> Edit\t\t\t"
	menu1 += "<f: cyan bold>F3:<t: white bold> Search\t\t\t<f: cyan bold>F4:<t: white bold> Delete\t\t\t<f: cyan bold>F5:<t: white bold> Exit"
	menu1 += strings.Repeat(" ", 80)
	menu1 += "<t: cyan bold>Using Database: <t: white bold>" + bookmarksPrefs.Database

	_ = ui.CreateLabel(mainWindow, winSize.windowWidth, 1, menu1, ui.Fixed)

	_ = ui.CreateLabel(mainWindow, winSize.screenWidth/2, 1, "", ui.Fixed)

	idNameDiv := strings.Repeat(" ", 10)
	nameTagsDiv := strings.Repeat(" ", 45)
	tagsUrlDiv := strings.Repeat(" ", 55)

	header := fmt.Sprintf("<t: white bold>ID%sNAME%sTAGS%sURL", idNameDiv, nameTagsDiv, tagsUrlDiv)

	_ = ui.CreateLabel(mainWindow, winSize.screenWidth/2, 1, header, ui.Fixed)

	listView := ui.CreateListBox(mainWindow, winSize.screenWidth-10, winSize.screeHeight-5, ui.Fixed)

	showWelcomeMessage(listView)

	ui.ActivateControl(mainWindow, listView)

	listView.OnKeyPress(func(key termbox.Key) bool {
		if key == termbox.KeyF1 {
			newBookmark(&winSize, bookmarksPrefs.Database)

		} else if key == termbox.KeyF2 {
			if listView.SelectedItemText() != "" {
				editBookmark(&winSize, bookmarksPrefs.Database, listView.SelectedItemText())
			}

		} else if key == termbox.KeyF3 {
			searchBookmark(&winSize, listView, bookmarksPrefs.Database)

		} else if key == termbox.KeyF4 {
			id := strings.Split(listView.SelectedItemText(), " ")[0]
			db.DeleteBookmark(id, bookmarksPrefs.Database)
			listviewActiveLine := listView.SelectedItem()
			listView.RemoveItem(listviewActiveLine)

		} else if key == termbox.KeyEnter {
			id := strings.Split(listView.SelectedItemText(), " ")[0]
			url := strings.Split(db.SearchByID(id, bookmarksPrefs.Database), "■")[2]
			cmd := exec.Command(bookmarksPrefs.Browser, url)
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}

		} else if key == termbox.KeyF5 {
			go ui.Stop()
		}
		return false
	})
}

func newBookmark(winSize *WindowSize, database string) {

	menu1 := "<f: cyan bold> F1:<t: white bold> Save\t\t\t<t: cyan bold>F2: <t: white bold>Clear Fields\t\t\t"
	menu1 += "<f: cyan bold>F3:<t: white bold> Cancel"

	mainWindow := ui.AddWindow(winSize.windowX, winSize.windowY, winSize.windowWidth, winSize.windowHeight, "")

	mainWindow.SetBorder(1)
	mainWindow.SetTitleButtons(0)
	mainWindow.SetTitle("<t: white bold> NEW ")

	mainWindow.SetPack(ui.Vertical)

	_ = ui.CreateLabel(mainWindow, 1, 1, menu1, ui.AutoSize)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameURL := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameURL.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameURL, 6, 1, " Url: ", ui.AutoSize)
	editFieldURL := ui.CreateEditField(frameURL, winSize.windowWidth-1, "", ui.AutoSize)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameName := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameName.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameName, 6, 1, "Name: ", ui.AutoSize)
	editFieldName := ui.CreateEditField(frameName, winSize.windowWidth-1, "", ui.AutoSize)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameTags := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameTags.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameTags, 6, 1, "Tags: ", ui.AutoSize)
	editFieldTags := ui.CreateEditField(frameTags, winSize.windowWidth-1, "", ui.AutoSize)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	lblStatus := ui.CreateLabel(mainWindow, 1, 1, " ", ui.AutoSize)

	ui.ActivateControl(mainWindow, editFieldURL)

	editFieldURL.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			returnValue := db.InsertIntoDB(editFieldName.Title(), editFieldTags.Title(), editFieldURL.Title(), database)
			lblStatus.SetTitle(returnValue)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")
			lblStatus.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})

	editFieldName.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			returnValue := db.InsertIntoDB(editFieldName.Title(), editFieldTags.Title(), editFieldURL.Title(), database)
			lblStatus.SetTitle(returnValue)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")
			lblStatus.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})

	editFieldTags.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			returnValue := db.InsertIntoDB(editFieldName.Title(), editFieldTags.Title(), editFieldURL.Title(), database)
			lblStatus.SetTitle(returnValue)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")
			lblStatus.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})
}

func editBookmark(winSize *WindowSize, database, listboxText string) {

	id := strings.Split(listboxText, "  ")[0]

	queryResult := db.SearchByID(id, database)

	name := strings.Split(queryResult, "■")[0]
	tags := strings.Split(queryResult, "■")[1]
	url := strings.Split(queryResult, "■")[2]

	menu1 := "<f: cyan bold> F1:<t: white bold> Search\t\t\t<t: cyan bold>F2: <t: white bold>Clear Fields\t\t\t"
	menu1 += "<f: cyan bold>F3:<t: white bold> Cancel"

	mainWindow := ui.AddWindow(winSize.windowX, winSize.windowY, winSize.windowWidth, winSize.windowHeight, "")

	mainWindow.SetBorder(1)
	mainWindow.SetTitleButtons(0)
	mainWindow.SetTitle("<t: white bold> EDIT ")

	mainWindow.SetPack(ui.Vertical)

	_ = ui.CreateLabel(mainWindow, 1, 1, menu1, ui.AutoSize)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameURL := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameURL.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameURL, 6, 1, " Url: ", ui.AutoSize)
	editFieldURL := ui.CreateEditField(frameURL, winSize.windowWidth-1, "", ui.AutoSize)
	editFieldURL.SetTitle(url)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameName := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameName.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameName, 6, 1, "Name: ", ui.AutoSize)
	editFieldName := ui.CreateEditField(frameName, winSize.windowWidth-1, "", ui.AutoSize)
	editFieldName.SetTitle(name)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameTags := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameTags.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameTags, 6, 1, "Tags: ", ui.AutoSize)
	editFieldTags := ui.CreateEditField(frameTags, winSize.windowWidth-1, "", ui.AutoSize)
	editFieldTags.SetTitle(tags)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	lblStatus := ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)

	ui.ActivateControl(mainWindow, editFieldURL)

	editFieldURL.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			returnValue := db.UpdateBookmark(id, editFieldName.Title(), editFieldTags.Title(), editFieldURL.Title(), database)
			lblStatus.SetTitle(returnValue)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")
			lblStatus.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})

	editFieldName.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			returnValue := db.UpdateBookmark(id, editFieldName.Title(), editFieldTags.Title(), editFieldURL.Title(), database)
			lblStatus.SetTitle(returnValue)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")
			lblStatus.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})

	editFieldTags.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			returnValue := db.UpdateBookmark(id, editFieldName.Title(), editFieldTags.Title(), editFieldURL.Title(), database)
			lblStatus.SetTitle(returnValue)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")
			lblStatus.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})
}

func searchBookmark(winSize *WindowSize, listbox *ui.ListBox, database string) {

	menu1 := "<f: cyan bold> F1:<t: white bold> Save\t\t\t<t: cyan bold>F2: <t: white bold>Clear\t\t\t"
	menu1 += "<f: cyan bold>F3:<t: white bold> Cancel"

	mainWindow := ui.AddWindow(winSize.windowX, winSize.windowY, winSize.windowWidth, winSize.windowHeight, "")

	mainWindow.SetBorder(1)
	mainWindow.SetTitleButtons(0)
	mainWindow.SetTitle("<t: white bold> SEARCH ")

	mainWindow.SetPack(ui.Vertical)

	_ = ui.CreateLabel(mainWindow, 1, 1, menu1, ui.AutoSize)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameURL := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameURL.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameURL, 9, 1, "  In url ", ui.AutoSize)
	editFieldURL := ui.CreateEditField(frameURL, winSize.windowWidth-1, "", ui.AutoSize)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameName := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameName.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameName, 9, 1, " In Name ", ui.AutoSize)
	editFieldName := ui.CreateEditField(frameName, winSize.windowWidth-1, "", ui.AutoSize)

	_ = ui.CreateLabel(mainWindow, 1, 1, "", ui.AutoSize)
	frameTags := ui.CreateFrame(mainWindow, winSize.windowWidth, 1, ui.BorderNone, ui.Fixed)
	frameTags.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frameTags, 9, 1, " In Tags ", ui.AutoSize)
	editFieldTags := ui.CreateEditField(frameTags, winSize.windowWidth-1, "", ui.AutoSize)

	ui.ActivateControl(mainWindow, editFieldURL)

	editFieldURL.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			listbox.Clear()
			db.SearchBookmark(listbox, "url", editFieldURL.Title(), database)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})

	editFieldName.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			db.SearchBookmark(listbox, "name", editFieldName.Title(), database)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})

	editFieldTags.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyF1 {
			db.SearchBookmark(listbox, "tags", editFieldTags.Title(), database)

		} else if key == termbox.KeyF2 {
			editFieldName.SetTitle("")
			editFieldTags.SetTitle("")
			editFieldURL.SetTitle("")

		} else if key == termbox.KeyF3 {
			ui.PutEvent(ui.Event{Type: ui.EventCloseWindow})
		}

		return false
	})
}

func showWelcomeMessage(listbox *ui.ListBox) {
	listbox.AddItem("")
	listbox.AddItem("██     ██ ███████ ██       ██████  ██████  ███    ███ ███████")
	listbox.AddItem("██     ██ ██      ██      ██      ██    ██ ████  ████ ██     ")
	listbox.AddItem("██  █  ██ █████   ██      ██      ██    ██ ██ ████ ██ █████  ")
	listbox.AddItem("██ ███ ██ ██      ██      ██      ██    ██ ██  ██  ██ ██     ")
	listbox.AddItem(" ███ ███  ███████ ███████  ██████  ██████  ██      ██ ███████")
	listbox.AddItem("")
}
