package main

import (
	"net/http"
	"io/ioutil"
	"html/template"
)

func logHandler(logFile string) func(res http.ResponseWriter, _ *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		logTemplate, _ := template.ParseFiles(templatePath("log.html"))
		logFileBytes, _ := ioutil.ReadFile(logFile)
		logTemplate.Execute(res, string(logFileBytes))
	}
}
