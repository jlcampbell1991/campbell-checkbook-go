package main

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title string
	Body  template.HTML
}

func main() {
	tmpl := template.Must(template.ParseFiles(
		"views/layouts/menu.html",
		"views/layouts/application.html",
	))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{Title: "Home", Body: template.HTML("<p>Welcome to Checkbook</p>")}
		if err := tmpl.ExecuteTemplate(w, "application", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
