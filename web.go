package main

import (
	"log"
	"net/http"
	"html/template"
	"os"
	"path"
	"strconv"
	"github.com/epond/porthole/music"
)

const latestAdditionsLimit  = 10

func main() {
	musicFolder := os.Getenv("MUSIC_FOLDER")
	knownReleasesFile := os.Getenv("KNOWN_RELEASES_FILE")
	knownReleasesBackup := os.Getenv("KNOWN_RELEASES_BACKUP")
	gitCommit := os.Getenv("GIT_COMMIT")
	logFile := os.Getenv("LOG_FILE")
	fetchInterval, _ := strconv.Atoi(os.Getenv("FETCH_INTERVAL"))
	dashboardRefreshInterval, _ := strconv.Atoi(os.Getenv("DASHBOARD_REFRESH_INTERVAL"))
	foldersToScan := os.Getenv("FOLDERS_TO_SCAN")

	log.Printf("Starting porthole. Music folder: %v, Known releases file: %v, Backup: %v", musicFolder, knownReleasesFile, knownReleasesBackup)

	recordCollectionAdditions := music.NewFileBasedAdditions(
		musicFolder,
		knownReleasesFile,
		knownReleasesBackup,
		foldersToScan,
		latestAdditionsLimit)
	statusCoordinator := NewStatusCoordinator(
		gitCommit,
		fetchInterval,
		recordCollectionAdditions)

	http.HandleFunc("/", templateHandler("dashboard.html", dashboardRefreshInterval * 1000))
	http.HandleFunc("/dashinfo", templateHandler("dashinfo.html", statusCoordinator.status))
	http.HandleFunc("/log", logHandler(logFile))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe(":9000", nil)
}

func templateHandler(templateFile string, data interface{}) func(res http.ResponseWriter, req *http.Request) {
	parsedTemplate, _ := template.ParseFiles(templatePath(templateFile))
	return func(res http.ResponseWriter, req *http.Request) {
		parsedTemplate.Execute(res, data)
	}
}

func templatePath(file string) string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/html/", file)
}