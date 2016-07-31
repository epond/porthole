package main

import (
	"log"
	"net/http"
	"html/template"
	"os"
	"path"
)

func main() {
	log.Print("Starting porthole...")

	status := &Status{0}

	NewStatusCoordinator(status, 2)

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/dashinfo", dashboardInfoHandler(status))

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe("localhost:9000", nil)
}

func dashboardHandler(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles(templatePath("dashboard.html"))
	t.Execute(res, nil)
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