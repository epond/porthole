package main

import (
	"net/http"
	"io/ioutil"
	"html/template"
	"fmt"
)

func logHandler(logFile string) func(res http.ResponseWriter, _ *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		logTemplate, _ := template.ParseFiles(templatePath("log.html"))
		logFileBytes, _ := ioutil.ReadFile(logFile)
		if len(logFileBytes) > 0 {
			logTemplate.Execute(res, string(logFileBytes))
		} else {
			logTemplate.Execute(res, fmt.Sprintf("there is no log at %v", logFile))
		}
	}
}
