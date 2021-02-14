package gui

import (
	"fmt"
	"geminiClientGo/coms"
	"log"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

const (
	defaultXAlign float64 = 0.1
	widgetHight           = 20
	urlInputWidth         = 500
	defaultURL            = "gemini://gemini.circumlunar.space/"
	prefixURL             = "gemini://"
	sufixURL              = "/"
	lineBreak             = "\n"
)

type App struct {
	mainWindow     gtk.Window
	mainBox        gtk.Box
	menuBox        gtk.Box
	scrolledWindow gtk.ScrolledWindow
	navBox         gtk.Box
}

var markupMap = map[int]string{
	1: "<span size=\"15000\"><b>%s</b></span>",
	2: "<span size=\"22000\"><b>%s</b></span>",
	3: "<span size=\"30000\"><b>%s</b></span>",
}

func CreateApp() {
	gtk.Init(nil)

	var app App = App{}
	app.configureWindow()

	gtk.Main()
}

func (app *App) configureWindow() {

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Uanble to create a new window: ", err)
	}

	app.mainWindow = *win

	app.mainWindow.SetTitle("Gemini Browser")
	app.mainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})

	app.mainWindow.SetDefaultSize(800, 600)

	app.addBoxes()

	app.mainWindow.ShowAll()
}

func (app *App) addBoxes() {
	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	if err != nil {
		log.Fatal("Error creating main box: ", err)
	}
	app.mainBox = *mainBox
	app.addUIComponents()
	app.mainWindow.Add(&app.mainBox)
}

func (app *App) addUIComponents() {
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
	urlInput.SetHExpand(true)
	urlInput.SetSizeRequest(urlInputWidth, widgetHight)

	sw, err := gtk.ScrolledWindowNew(
		app.mainWindow.GetFocusHAdjustment(), app.mainWindow.GetFocusVAdjustment())
	if err != nil {
		log.Fatal("Error creating scrolled window: ", err)
	}
	sw.SetVExpand(true)

	app.scrolledWindow = *sw

	navLblBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	if err != nil {
		log.Fatal("Error creating page box: ", err)
	}
	navLblBox.SetVExpand(false)
	app.navBox = *navLblBox
	app.scrolledWindow.Add(&app.navBox)

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
		app.checkURL(url)
	})
	navButtonBox.Add(back)
	navButtonBox.Add(ford)
	navButtonBox.Add(urlInput)
	navButtonBox.Add(nav)
	app.mainBox.Add(navButtonBox)
	app.mainBox.Add(&app.scrolledWindow)
}

func (app *App) checkURL(url string) {
	if url == "" {
		return
	}
	if !strings.HasPrefix(url, prefixURL) {
		url = prefixURL + url
	}
	if !strings.HasSuffix(url, sufixURL) {
		url = url + sufixURL
	}
	app.navigate(url)
}

func (app *App) navigate(url string) {
	response := coms.Request(url)

	app.destoryNavBoxChildren()
	app.formatResponse(response)
	// app.navBox.Add(lbl)
	app.mainWindow.ShowAll()
}

func (app *App) destoryNavBoxChildren() {
	app.navBox.GetChildren().Foreach(func(item interface{}) {
		item.(*gtk.Widget).Destroy()
	})
}

func (app *App) formatResponse(response string) {
	for _, line := range strings.Split(response, lineBreak) {
		if strings.HasPrefix(line, "```") {
			continue
		} else if strings.HasPrefix(line, "###") {
			app.addLabel(line, 3)
		} else if strings.HasPrefix(line, "##") {
			app.addLabel(line, 2)
		} else if strings.HasPrefix(line, "#") {
			app.addLabel(line, 1)
		} else {
			app.addLabel(line, 0)
		}
	}
}

func (app *App) addLabel(text string, markup int) {
	lbl, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Error creatign label for page: ", err)
	}
	lbl.SetLineWrap(true)
	lbl.SetVExpand(false)

	labAlign := getLabAlign(markup)
	lbl.SetXAlign(labAlign)

	if markup == 0 {
		lbl.SetText(text)
	} else {
		lbl.SetMarkup(fmt.Sprintf(markupMap[1], text))
	}

	app.navBox.Add(lbl)
}

func getLabAlign(markup int) float64 {
	if markup == 0 {
		return defaultXAlign
	}

	fmt.Println(defaultXAlign - (float64(markup) / 100))
	return defaultXAlign - (float64(markup) / 100)
}
