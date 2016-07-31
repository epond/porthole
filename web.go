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

	http.HandleFunc("/", DashboardHandler(status))

	log.Print("porthole active - browse to http://localhost:9000")
	http.ListenAndServe("localhost:9000", nil)
}

func DashboardHandler(status *Status) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Counter:%v (E)\n", status.Counter)
	}
}
