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
)

func main() {
	config := hub.ConfigFromEnv()

	log.Printf("Starting porthole. Music folder: %v, Known albums file: %v, Backup: %v, Folders to scan: %v", config.MusicFolder, config.KnownAlbumsFile, config.KnownAlbumsBackup, config.FoldersToScan)

	scanner := &foldermusic.DepthAwareFolderScanner{
		config.MusicFolder,
		config.FoldersToScan,
	}
	persister := &foldermusic.KnownAlbumsWithBackup{
		config.KnownAlbumsFile,
		config.KnownAlbumsBackup,
	}
	clock := NewClockTicker(time.Duration(config.FetchInterval) * time.Millisecond)

	porthole := hub.NewPorthole(
		scanner,
		persister,
		clock,
		config,
	)

	http.HandleFunc("/", templateHandler("dashboard.html", config.DashboardRefreshInterval))
	http.HandleFunc("/dashinfo", templateHandlerDynamic("dashinfo.html", func() interface{} {
		return porthole.GetStatus()
	}))
	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		porthole.RequestScan()
	})
	http.HandleFunc("/log", logHandler(config.LogFile))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe(":9000", nil)
}

func templateHandler(templateFile string, data interface{}) func(res http.ResponseWriter, req *http.Request) {
	return templateHandlerDynamic(templateFile, func() interface{} { return data })
}

func templateHandlerDynamic(templateFile string, data func() interface{}) func(res http.ResponseWriter, req *http.Request) {
	parsedTemplate, _ := template.ParseFiles(templatePath(templateFile))
	return func(res http.ResponseWriter, req *http.Request) {
		parsedTemplate.Execute(res, data())
	}
}

func templatePath(file string) string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/html/", file)
}
