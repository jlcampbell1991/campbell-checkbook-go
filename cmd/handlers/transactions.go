package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// TransactionsIndex handles GET /transactions, renders a paginated list.
func TransactionsIndex(w http.ResponseWriter, r *http.Request) {
	// Parse page query param, default 0.
	page := 0
	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p >= 0 {
			page = p
		}
	}

	// Retrieve transactions (stub repository).
	txs, err := db.Index(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build next URL.
	nextURL := fmt.Sprintf("/transactions?page=%d", page+1)

	// Parse templates.
	tmpl, err := template.ParseGlob("views/transactions/*.tmpl")
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Data passed to template.
	data := struct {
		Transactions []domain.Transaction
		NextURL      string
		Now          time.Time
	}{
		Transactions: txs,
		NextURL:      nextURL,
		Now:          time.Now(),
	}

	// Execute the layout within the application template if needed.
	// Assuming the transaction templates include "application.tmpl" via inheritance.
	if err := tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
