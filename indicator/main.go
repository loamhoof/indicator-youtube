package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	pb "github.com/loamhoof/indicator"
	"github.com/loamhoof/indicator/client"
)

const (
	ID = "indicator-youtube"
)

var (
	play, pause, logFile string
	port                 int
	sc                   *client.ShepherdClient
	playing              string
	logger               *log.Logger
	resetTimer           *time.Timer
	resetAfter           time.Duration = time.Second * 3
)

func init() {
	flag.IntVar(&port, "port", 15000, "Port of the shepherd")
	flag.StringVar(&play, "play", "", "Path to the play icon")
	flag.StringVar(&pause, "pause", "", "Path to the pause icon")
	flag.StringVar(&logFile, "log", "", "Log file")

	flag.Parse()

	logger = log.New(os.Stdout, "", log.LstdFlags)
}

func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger.Println("Request", req.URL.Path)

	resetTimer.Reset(resetAfter)

	path := decodePath(req.URL)

	title := cleanTitle(path[1])
	current := formatDuration(path[2])
	duration := formatDuration(path[3])
	label := fmt.Sprintf("%s (%s / %s)", title, current, duration)

	var icon string
	if path[4] == "false" {
		icon = play
	} else {
		icon = pause
	}

	iReq := &pb.Request{
		Id:     ID,
		Label:  label,
		Icon:   icon,
		Active: true,
	}
	if _, err := sc.Update(iReq); err != nil {
		logger.Println(err)
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

func main() {
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			logger.Fatalln(err)
		}
		defer f.Close()
		logger = log.New(f, "", log.LstdFlags)
	}

	sc = client.NewShepherdClient(port)
	for {
		err := sc.Init()
		if err == nil {
			break
		}
		logger.Fatalf("Could not connect: %v", err)

		time.Sleep(time.Second * 5)
	}
	defer sc.Close()

	resetTimer = time.AfterFunc(resetAfter, func() {
		iReq := &pb.Request{
			Id:     ID,
			Active: false,
		}
		if _, err := sc.Update(iReq); err != nil {
			logger.Println(err)
		}
	})

	serve()
}
