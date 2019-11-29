package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/epond/porthole/foldermusic"
	"github.com/epond/porthole/hub"
	"github.com/epond/porthole/status"
)

func main() {
	config := hub.ConfigFromEnv()

	log.Printf("Starting porthole. Music folder: %v, Known albums file: %v, Backup: %v, Folders to scan: %v", config.MusicFolder, config.KnownAlbumsFile, config.KnownAlbumsBackup, config.FoldersToScan)

	folderScanner := &foldermusic.DepthAwareFolderScanner{
		config.MusicFolder,
		config.FoldersToScan,
	}
	knownAlbums := &foldermusic.KnownAlbumsWithBackup{
		config.KnownAlbumsFile,
		config.KnownAlbumsBackup,
	}
	albumAdditions := foldermusic.NewAdditions(
		folderScanner,
		knownAlbums,
		config.LatestAdditionsLimit)
	clock := time.Tick(time.Duration(config.FetchInterval) * time.Millisecond)
	statusCoordinator := status.NewCoordinator(
		config.GitCommit,
		&status.MusicStatusWorker{albumAdditions},
		clock,
		time.Duration(config.SleepAfter)*time.Millisecond)

	http.HandleFunc("/", templateHandler("dashboard.html", config.DashboardRefreshInterval))
	http.HandleFunc("/dashinfo", dashinfoHandler(statusCoordinator.Status))
	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		statusCoordinator.Status.LastRequest = time.Now()
	})
	http.HandleFunc("/log", logHandler(config.LogFile))
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
