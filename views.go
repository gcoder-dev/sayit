package main

import (
	"html/template"
	"log"
	"net/http"
)

func landingView(writer http.ResponseWriter, request *http.Request) {
	temp, err := template.ParseFiles("landing.html")
	if err != nil {
		log.Println(err)
		return
	}
	err = temp.Execute(writer, nil)
	if err != nil {
		return
	}
}
