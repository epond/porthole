package main

import (
	"log"
	"net/http"
	"html/template"
	"os"
	"path"
	"fmt"
	"io/ioutil"
)

func main() {
	musicFolder := os.Getenv("MUSIC_FOLDER")
	knownReleasesFile := os.Getenv("KNOWN_RELEASES_FILE")
	gitCommit := os.Getenv("GIT_COMMIT")
	logFile := os.Getenv("LOG_FILE")
	status := &Status{
		GitCommit: gitCommit,
		Counter: 0,
		LatestAdditions: []string{},
	}

	log.Printf("Starting porthole. Music folder: %v, Known releases file: %v", musicFolder, knownReleasesFile)

	NewStatusCoordinator(status, musicFolder, knownReleasesFile, 30)

	http.HandleFunc("/", dashboardHandler())
	http.HandleFunc("/dashinfo", dashboardInfoHandler(status))
	http.HandleFunc("/log", logHandler(logFile))

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe("localhost:9000", nil)
}

func dashboardHandler() func(res http.ResponseWriter, req *http.Request) {
	dashboard, _ := ioutil.ReadFile(templatePath("dashboard.html"))
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, string(dashboard))
	}
}

func dashboardInfoHandler(status *Status) func(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles(templatePath("dashinfo.html"))
	return func(res http.ResponseWriter, req *http.Request) {
		t.Execute(res, status)
	}
}

func templatePath(file string) string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/html/", file)
}