package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/conformal/gotk3/gtk"
	"github.com/doxxan/appindicator"
	"github.com/doxxan/appindicator/gtk-extensions/gotk3"
)

var (
	icon      string
	indicator *gotk3.AppIndicatorGotk3
)

func init() {
	flag.StringVar(&icon, "icon", "", "Path to the icon")

	flag.Parse()
}

func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		indicator.SetLabel("∅", "")

		return
	}

	title := req.URL.Path[1 : len(req.URL.Path)-10]
	indicator.SetLabel(title, "")
}

func serve() {
	http.HandleFunc("/", ServeHTTP)
	log.Println("Listening...")
	if err := http.ListenAndServe(":12346", nil); err != nil {
		log.Println(err)
	}
}

func indicate() {
	gtk.Init(nil)

	indicator = gotk3.NewAppIndicator("indicator-youtube", icon, appindicator.CategorySystemServices)

	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetLabel("∅", "")

	menu, err := gtk.MenuNew()
	if err != nil {
		panic(err)
	}

	menuItem, err := gtk.MenuItemNewWithLabel("")
	if err != nil {
		panic(err)
	}

	menu.Append(menuItem)

	menuItem.Show()
	indicator.SetMenu(menu)

	gtk.Main()
}

func main() {
	go serve()

	indicate()
}
