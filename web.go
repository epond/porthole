package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting porthole...")

	NewStatusCoordinator(2)

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/dashboard", DashboardHandler)

	log.Print("porthole active - go to http://localhost:9000/dashboard")
	http.ListenAndServe("localhost:9000", nil)
}

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "GET /dashboard (see project readme for more information)")
}

func DashboardHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Dashboard will appear here")
}
