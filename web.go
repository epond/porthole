package main

import (
	"log"
	"net/http"
	"html/template"
	"os"
	"path"
	"strconv"
)

func main() {
	musicFolder := os.Getenv("MUSIC_FOLDER")
	knownReleasesFile := os.Getenv("KNOWN_RELEASES_FILE")
	gitCommit := os.Getenv("GIT_COMMIT")
	logFile := os.Getenv("LOG_FILE")
	fetchInterval, _ := strconv.Atoi(os.Getenv("FETCH_INTERVAL"))
	dashboardRefreshInterval, _ := strconv.Atoi(os.Getenv("DASHBOARD_REFRESH_INTERVAL"))
	status := &Status{
		GitCommit: gitCommit,
		Counter: 0,
		LatestAdditions: []string{},
	}

	log.Printf("Starting porthole. Music folder: %v, Known releases file: %v", musicFolder, knownReleasesFile)

	NewStatusCoordinator(status, musicFolder, knownReleasesFile, fetchInterval)

	http.HandleFunc("/", dashboardHandler(dashboardRefreshInterval * 1000))
	http.HandleFunc("/dashinfo", dashboardInfoHandler(status))
	http.HandleFunc("/log", logHandler(logFile))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe("localhost:9000", nil)
}

func dashboardHandler(dashboardRefreshInterval int) func(res http.ResponseWriter, req *http.Request) {
	dashboardTemplate, _ := template.ParseFiles(templatePath("dashboard.html"))
	return func(res http.ResponseWriter, req *http.Request) {
		dashboardTemplate.Execute(res, dashboardRefreshInterval)
	}
}

func dashboardInfoHandler(status *Status) func(res http.ResponseWriter, req *http.Request) {
	dashinfoTemplate, _ := template.ParseFiles(templatePath("dashinfo.html"))
	return func(res http.ResponseWriter, req *http.Request) {
		dashinfoTemplate.Execute(res, status)
	}
}

func templatePath(file string) string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/html/", file)
}