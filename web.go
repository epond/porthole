package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting porthole... go to http://localhost:9000/dashboard")
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/dashboard", DashboardHandler)

	http.ListenAndServe("localhost:9000", nil)
}

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "GET /dashboard (see project readme for more information)")
}

func DashboardHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Dashboard will appear here")
}
