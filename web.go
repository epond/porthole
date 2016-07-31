package main

import (
	"log"
	"net/http"
	"html/template"
	"os"
	"path"
	"fmt"
)

func main() {
	log.Print("Starting porthole...")

	status := &Status{0}

	NewStatusCoordinator(status, 2)

	http.HandleFunc("/", DashboardHandler(status))
	http.HandleFunc("/dashinfo", DashboardInfoHandler(status))

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe("localhost:9000", nil)
}

func DashboardHandler(status *Status) func(res http.ResponseWriter, req *http.Request) {
	templatePath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole", "html/dashboard.html")
	log.Printf("Loading dashboard template from %v", templatePath)
	t, _ := template.ParseFiles(templatePath)
	return func(res http.ResponseWriter, req *http.Request) {
		t.Execute(res, status)
	}
}

func DashboardInfoHandler(status *Status) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Counter:%v\n", status.Counter)
	}
}