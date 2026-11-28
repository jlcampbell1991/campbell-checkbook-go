package main

import (
	"embed"
	"html/template"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/config"
	"github.com/jlcampbell1991/campbell-checkbook-go/cmd/handlers"
)

// Template helper functions
func IsWithdrawal(amount float64) bool { return amount < 0 }

func CurrencyFromDouble(v float64) string { return fmt.Sprintf("%.2f", v) }

//go:embed views/layouts/*.tmpl views/layouts/*.html views/transactions/*.tmpl
var templateFS embed.FS

func main() {
	// Parse all layout templates
	tmpl, err := template.ParseFS(templateFS, "views/layouts/*.tmpl", "views/layouts/*.html", "views/transactions/*.tmpl")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	// Register template functions
	tmpl = tmpl.Funcs(template.FuncMap{"IsWithdrawal": IsWithdrawal, "CurrencyFromDouble": CurrencyFromDouble})

	config.Load()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static assets if present (placeholder paths)
	staticDir := http.Dir(filepath.Join("static"))
	r.Handle("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(staticDir)))
	r.Handle("/javascripts/*", http.StripPrefix("/javascripts/", http.FileServer(staticDir)))

	// Register line‑item routes
	r.Mount("/line-items", handlers.LineItemRoutes(tmpl, nil)) // TODO: provide a concrete LineItemRepository implementation
	r.Get("/transactions", handlers.TransactionsIndex)

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
