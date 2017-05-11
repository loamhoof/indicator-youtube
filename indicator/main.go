package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/conformal/gotk3/gtk"
	"github.com/doxxan/appindicator"
	"github.com/doxxan/appindicator/gtk-extensions/gotk3"
)

var (
	play, pause    string
	state, playing string
	indicator      *gotk3.AppIndicatorGotk3
)

func init() {
	flag.StringVar(&play, "play", "", "Path to the play icon")
	flag.StringVar(&pause, "pause", "", "Path to the pause icon")

	flag.Parse()
}

func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if indicator == nil {
		return
	}

	path := decodePath(req.URL)

	oldState := state
	state := path[0]

	switch state {
	case "play":
		title := cleanTitle(path[1])
		playing = title

		current := formatDuration(path[2])
		duration := formatDuration(path[3])

		if state != oldState {
			indicator.SetIcon(play, "")
		}
		indicator.SetLabel(fmt.Sprintf("%s (%s / %s)", title, current, duration), "")

		if indicator.GetStatus() != appindicator.StatusActive {
			indicator.SetStatus(appindicator.StatusActive)
		}
	case "pause":
		for i := 1; i < len(path); i++ {
			title := cleanTitle(path[i])
			if title == playing {
				if state != oldState {
					indicator.SetIcon(pause, "")
				}

				return
			}
		}

		if indicator.GetStatus() != appindicator.StatusPassive {
			indicator.SetStatus(appindicator.StatusPassive)
		}
	default:
	}
}

func decodePath(u *url.URL) []string {
	if u.RawPath == "" {
		return strings.Split(u.Path[1:], "/")
	}

	encodedPath := strings.Split(u.RawPath[1:], "/")

	decodedPath := make([]string, len(encodedPath))
	for i, encodedComponent := range encodedPath {
		decodedComponent, _ := url.PathUnescape(encodedComponent)
		decodedPath[i] = decodedComponent
	}

	return decodedPath
}

func cleanTitle(title string) string {
	if len(title) < 10 {
		return ""
	}

	return title[:len(title)-10]
}

func formatDuration(seconds string) string {
	d, _ := time.ParseDuration(seconds + "s")

	h := int(d / time.Hour)
	m := int(d/time.Minute) % 60
	s := int(d/time.Second) % 60

	if h == 0 {
		return fmt.Sprintf("%v:%02v", m, s)
	}

	return fmt.Sprintf("%v:%v:%02v", h, m, s)
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

	indicator = gotk3.NewAppIndicator("indicator-youtube", pause, appindicator.CategorySystemServices)

	indicator.SetStatus(appindicator.StatusPassive)

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
