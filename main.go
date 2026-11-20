package main

import (
	"embed"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/config"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jlcampbell1991/campbell-checkbook-go/cmd/handlers"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"github.com/jlcampbell1991/campbell-checkbook-go/utils"
	"fmt"
	"strconv"
)

//go:embed views/layouts/*.tmpl views/layouts/*.html views/transactions/*.tmpl views/line_items/*.tmpl
var templateFS embed.FS

func main() {
	// Parse all layout templates
	tmpl, err := template.ParseFS(templateFS, "views/layouts/*.tmpl", "views/layouts/*.html", "views/transactions/*.tmpl", "views/line_items/*.tmpl")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}
	tmpl = tmpl.Funcs(template.FuncMap{"isWithdrawal": utils.IsWithdrawal})

	// Make templates available to transaction handler
	handlers.SetTemplate(tmpl)

	config.Load()

	r := chi.NewRouter()

	// Initialize repository (stub implementation)
	repo := db.NewStubAccountRepo()
	lineItemRepo := db.NewStubLineItemRepo()
	transactionRepo := db.NewStubTransactionRepo()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static assets if present (placeholder paths)
	staticDir := http.Dir(filepath.Join("static"))
	r.Handle("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(staticDir)))
	r.Handle("/javascripts/*", http.StripPrefix("/javascripts/", http.FileServer(staticDir)))

	// New line-item creation routes
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
		accountId := r.FormValue("accountId")
		if name == "" {
			http.Error(w, "name required", http.StatusBadRequest)
			return
		}
		bal, err := strconv.ParseFloat(balStr, 64)
		if err != nil {
			http.Error(w, "invalid balance", http.StatusBadRequest)
			return
		}
		li := domain.LineItem{Name: name, Balance: bal, AccountID: accountId}
		if _, err := lineItemRepo.Insert(r.Context(), li); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/accounts/%s", accountId), http.StatusSeeOther)
	})

	// Mount account routes
	r.Mount("/accounts", handlers.AccountRoutes(tmpl, repo))
	r.Mount("/line-items", handlers.LineItemRoutes(tmpl, lineItemRepo))

	// Transactions page
	r.Mount("/transactions", handlers.TransactionRoutes(tmpl, transactionRepo))

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
