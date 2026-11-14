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
		filepath.Join("views", "line_items", "add.html"),
		filepath.Join("views", "transactions", "index.html"),
		filepath.Join("views", "transfers", "form.html"),
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

	// Line items: show form to add new line item (GET)
	r.Get("/accounts/{id}/line-items/new", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		data := struct {
			Title     string
			AccountID string
		}{Title: "New Line Item", AccountID: id}
		if err := tmpl.ExecuteTemplate(w, "lineitem_add", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Handle line item creation (POST)
	r.Post("/accounts/{id}/line-items", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}
		id := chi.URLParam(r, "id")
		name := r.FormValue("name")
		balance := r.FormValue("balance")
		_ = name
		_ = balance
		// TODO: persist line item
		_ = id // suppress unused warning
		http.Redirect(w, r, "/accounts/"+id, http.StatusSeeOther)
	})

	// Transactions index page
	r.Get("/transactions", func(w http.ResponseWriter, r *http.Request) {
		data := struct{ Title string }{Title: "Transactions"}
		if err := tmpl.ExecuteTemplate(w, "transactions_index", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Transfers: show form (GET)
	r.Get("/transfers/new", func(w http.ResponseWriter, r *http.Request) {
		lineItemID := r.URL.Query().Get("line-item-id")
		// Mock line items data; in real app fetch from DB.
		type Item struct {
			Name    string
			ID      string
			Balance float64
		}
		lineItems := []Item{{Name: "Sample Item", ID: "1", Balance: 100.0}}
		data := struct {
			Title      string
			LineItemID string
			LineItems  []Item
		}{
			Title:      "New Transfer",
			LineItemID: lineItemID,
			LineItems:  lineItems,
		}
		if err := tmpl.ExecuteTemplate(w, "transfer_form", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Transfers: handle submission (POST)
	r.Post("/transfers", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}
		_ = r.FormValue("memo")
		_ = r.FormValue("amount")
		_ = r.FormValue("from")
		_ = r.FormValue("to")
		// TODO: persist transfer
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Start server
	addr := ":" + portStr
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
