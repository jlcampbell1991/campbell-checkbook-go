package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jlcampbell1991/campbell-checkbook-go/cmd/handlers"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/config"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// Template helper functions
func IsWithdrawal(amount float64) bool { return amount < 0 }

func CurrencyFromDouble(v float64) string { return fmt.Sprintf("%.2f", v) }

//go:embed views/layouts/*.tmpl views/layouts/*.html views/transactions/*.tmpl
var templateFS embed.FS

func main() {
	// Parse all layout templates
	tmpl, err := template.ParseFS(templateFS, "views/layouts/*.tmpl", "views/layouts/*.html", "views/transactions/*.tmpl", "views/line_items/*.tmpl")
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

	// New Line Item form routes
	r.Get("/accounts/{id}/line-items/new", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		data := struct {
			Title     string
			Action    string
			Method    string
			AccountID string
		}{
			Title:     "New Line Item",
			Action:    fmt.Sprintf("/accounts/%s/line-items", id),
			Method:    "post",
			AccountID: id,
		}
		if err := tmpl.ExecuteTemplate(w, "line_items_form.tmpl", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	r.Post("/accounts/{id}/line-items", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		balStr := r.FormValue("balance")
		acct := r.FormValue("accountId")
		if name == "" {
			http.Error(w, "missing name", http.StatusBadRequest)
			return
		}
		bal, err := strconv.ParseFloat(balStr, 64)
		if err != nil {
			http.Error(w, "invalid balance", http.StatusBadRequest)
			return
		}
		li := domain.LineItem{Name: name, Balance: bal, AccountID: acct}
		if _, err := db.LineItemRepo.Insert(r.Context(), li); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/accounts/%s", acct), http.StatusSeeOther)
	})

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
