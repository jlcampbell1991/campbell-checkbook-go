package main

import (
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Load environment variables from .env if present
	_ = godotenv.Load()
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	// Validate port number
	if _, err := strconv.Atoi(portStr); err != nil {
		log.Fatalf("Invalid PORT: %v", err)
	}

	// Serve static assets
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Render menu on root path
	tmpl := template.Must(template.ParseFiles("templates/menu.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Printf("Server starting on :%s", portStr)
	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
