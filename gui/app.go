package gui

import (
	"geminiClientGo/coms"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

const (
	widgetHight   = 20
	urlInputWidth = 500
	defaultURL    = "gemini://gemini.circumlunar.space/"
)

func CreateApp() {
	gtk.Init(nil)

	configureWindow()

	gtk.Main()
}

func configureWindow() gtk.Window {

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Uanble to create a new window: ", err)
	}

	win.SetTitle("Gemini Browser")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.SetDefaultSize(800, 600)

	addBoxes(win)

	win.ShowAll()

	return *win
}

func addBoxes(win *gtk.Window) {
	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	if err != nil {
		log.Fatal("Error creating main box: ", err)
	}
	addMenuButtons(mainBox)
	win.Add(mainBox)
}

func addMenuButtons(win *gtk.Box) {
	navButtonBox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		log.Fatal("Error creating navigation button box: ", err)
	}
	navButtonBox.SetHAlign(1)
	navButtonBox.SetVAlign(1)

	back, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal("Error creating navigate button: ", err)
	}
	back.SetLabel("BACK")

	ford, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal("Error creating navigate button: ", err)
	}
	ford.SetLabel("FORD")

	urlInput, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Error creating URL input: ", err)
	}
	urlInput.SetSizeRequest(urlInputWidth, widgetHight)

	sw, err := gtk.ScrolledWindowNew(win.GetFocusHAdjustment(), win.GetFocusVAdjustment())
	sw.SetVExpand(true)
	navLblBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	if err != nil {
		log.Fatal("Error creating page box: ", err)
	}
	navLblBox.SetHExpand(true)

	sw.Add(navLblBox)

	nav, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal("Error creating navigate button: ", err)
	}
	nav.SetLabel("NAV")
	nav.Connect("clicked", func() {
		url, err := urlInput.GetText()
		if err != nil {
			url = defaultURL
		}
		navigate(win, url, navLblBox)
	})

	navButtonBox.Add(back)
	navButtonBox.Add(ford)
	navButtonBox.Add(urlInput)
	navButtonBox.Add(nav)
	win.Add(navButtonBox)
	win.Add(sw)
}

func navigate(win *gtk.Box, url string, navLblBox *gtk.Box) {
	response := coms.Request(url)
	lbl, err := gtk.LabelNew(response)
	if err != nil {
		log.Fatal("Error creatign label for page: ", err)
	}
	lbl.SetLineWrap(true)
	lbl.SetWidthChars(30)

	navLblBox.Add(lbl)
	win.ShowAll()
}
