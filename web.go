package main

import (
	"log"
	"net/http"
	"html/template"
)

func main() {
	log.Print("Starting porthole...")

	status := &Status{0}

	NewStatusCoordinator(status, 2)

	http.HandleFunc("/", DashboardHandler(status))

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe("localhost:9000", nil)
}

func DashboardHandler(status *Status) func(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("html/dashboard.html")
	return func(res http.ResponseWriter, req *http.Request) {
		t.Execute(res, status)
	}
}
