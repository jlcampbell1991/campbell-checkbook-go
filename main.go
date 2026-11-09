package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func init() {
	// Parse layout and menu templates
	tmpl = template.Must(template.ParseFiles(
		"templates/layouts/application.html",
		"templates/layouts/menu.html",
	))
}

type pageData struct {
	Title string
	Body  template.HTML
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		Title: "Checkbook",
		Body:  "Welcome to Checkbook",
	}
	if err := tmpl.ExecuteTemplate(w, "application", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	// Serve static files from ./static if present
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := ":8080"
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
