package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/epond/porthole/foldermusic"
	"github.com/epond/porthole/status"
)

func main() {
	musicFolder := os.Getenv("MUSIC_FOLDER")
	knownReleasesFile := os.Getenv("KNOWN_RELEASES_FILE")
	knownReleasesBackup := os.Getenv("KNOWN_RELEASES_BACKUP")
	gitCommit := os.Getenv("GIT_COMMIT")
	logFile := os.Getenv("LOG_FILE")
	fetchInterval, _ := strconv.Atoi(os.Getenv("FETCH_INTERVAL"))
	dashboardRefreshInterval, _ := strconv.Atoi(os.Getenv("DASHBOARD_REFRESH_INTERVAL"))
	sleepAfter, _ := strconv.Atoi(os.Getenv("SLEEP_AFTER"))
	foldersToScan := os.Getenv("FOLDERS_TO_SCAN")
	latestAdditionsLimit, _ := strconv.Atoi(os.Getenv("LATEST_ADDITIONS_LIMIT"))

	log.Printf("Starting porthole. Music folder: %v, Known releases file: %v, Backup: %v, Folders to scan: %v", musicFolder, knownReleasesFile, knownReleasesBackup, foldersToScan)

	recordCollectionAdditions := foldermusic.NewAdditions(
		musicFolder,
		knownReleasesFile,
		knownReleasesBackup,
		foldersToScan,
		latestAdditionsLimit)
	clock := time.Tick(time.Duration(fetchInterval) * time.Millisecond)
	statusCoordinator := status.NewCoordinator(
		gitCommit,
		&status.MusicStatusWorker{recordCollectionAdditions},
		clock,
		time.Duration(sleepAfter)*time.Millisecond)

	http.HandleFunc("/", templateHandler("dashboard.html", dashboardRefreshInterval))
	http.HandleFunc("/dashinfo", dashinfoHandler(statusCoordinator.Status))
	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		statusCoordinator.Status.LastRequest = time.Now()
	})
	http.HandleFunc("/log", logHandler(logFile))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe(":9000", nil)
}

func dashinfoHandler(status *status.Status) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		templateHandler("dashinfo.html", status)(res, req)
	}
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
