package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

//go:embed views/layouts/*.tmpl views/layouts/*.html
var templateFS embed.FS

func main() {
	// Parse all layout templates
	tmpl, err := template.ParseFS(templateFS, "views/layouts/*.tmpl", "views/layouts/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	// Load environment variables and run DB migrations
	_ = godotenv.Load()
	cfg := db.PostgresConfig{
		User: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host: os.Getenv("POSTGRES_HOST"),
		DB: os.Getenv("POSTGRES_DB"),
		Port: os.Getenv("POSTGRES_PORT"),
	}
	dsn := db.BuildPostgresURL(cfg)
	db.RunMigrations(dsn)

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
	log.Printf("Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
