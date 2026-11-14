package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Determine port from environment variable, default 8080
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	// Validate port is numeric
	if _, err := strconv.Atoi(portStr); err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}

	// Initialize router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static assets
	r.Handle("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("static/stylesheets"))))
	r.Handle("/javascripts/*", http.StripPrefix("/javascripts/", http.FileServer(http.Dir("static/javascripts"))))

	// Parse templates
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "application.html"),
		filepath.Join("views", "layouts", "menu.html"),
		filepath.Join("views", "accounts", "add.html"),
	))

	// Simple home handler demonstrating layout usage
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct{ Title string }{Title: "Home"}
		if err := tmpl.ExecuteTemplate(w, "application", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Account creation form (GET)
	r.Get("/accounts/new", func(w http.ResponseWriter, r *http.Request) {
		data := struct{ Title string }{Title: "New Account"}
		if err := tmpl.ExecuteTemplate(w, "add", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Handle form submission (POST)
	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		// In a full implementation, we would persist the account.
		// For now, just redirect to home.
		_ = name // suppress unused warning
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Start server
	addr := ":" + portStr
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
