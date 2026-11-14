package main

import (
<<<<<<< HEAD
    "log"
    "net/http"
    "os"
    "strconv"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/joho/godotenv"
    "html/template"
    "path/filepath"
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
    // staticDir := http.Dir("static") // retained for reference
    r.Handle("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir(filepath.Join("static", "stylesheets")))) )
    r.Handle("/javascripts/*", http.StripPrefix("/javascripts/", http.FileServer(http.Dir(filepath.Join("static", "javascripts")))) )

    // Parse templates
    tmpl := template.Must(template.ParseFiles(
        filepath.Join("views", "layouts", "application.html"),
        filepath.Join("views", "layouts", "menu.html"),
    ))

    // Simple home handler demonstrating layout usage
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        data := struct{ Title string }{Title: "Home"}
        if err := tmpl.ExecuteTemplate(w, "application", data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })

    // Start server
    addr := ":" + portStr
    log.Printf("Starting server on %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
=======
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
>>>>>>> origin/task/add-nav-and-menu-7
}

