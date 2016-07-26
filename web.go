package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting porthole...")

	status := &Status{0}

	NewStatusCoordinator(status, 2)

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/dashboard", DashboardHandler(status))

	log.Print("porthole active - go to http://localhost:9000/dashboard")
	http.ListenAndServe("localhost:9000", nil)
}

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "GET /dashboard (see project readme for more information)")
}

func DashboardHandler(status *Status) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Counter:%v\n", status.Counter)
	}
}