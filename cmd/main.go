package main

import (
	"database/sql"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/config"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

//go:embed views/layouts/*.tmpl views/layouts/*.html
var templateFS embed.FS

func main() {
	// Load DB configuration from environment (and .env file if present)
	cfg := db.Load()
	dsn := db.BuildPostgresURL(cfg)

	// Test DB connection
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open DB connection:", err)
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		log.Fatal("failed to ping DB:", err)
	}
	log.Println("DB connection successful")

	// Run migrations (no-op if none present)
	db.RunMigrations(dsn)

	// Parse all layout templates
	tmpl, err := template.ParseFS(templateFS, "views/layouts/*.tmpl", "views/layouts/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	config.Load()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static assets if present (placeholder paths)
	staticDir := http.Dir(filepath.Join("static"))
	r.Handle("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(staticDir)))
	r.Handle("/javascripts/*", http.StripPrefix("/javascripts/", http.FileServer(staticDir)))

	// Home page handler
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title string
			Body  template.HTML
		}{
			Title: "Home",
			Body:  template.HTML("<p>Welcome to Campbell Checkbook.</p>"),
		}
		// Execute the application layout template
		if err := tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
