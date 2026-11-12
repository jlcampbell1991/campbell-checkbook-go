package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook/internal/domain"
)

// RegisterTransactions registers transaction routes on the provided router.
func RegisterTransactions(r chi.Router, lineItemRepo *db.InMemoryLineItemRepo) {
	// List transactions (placeholder implementation)
	r.Get("/transactions", listTransactions())
	// Transaction creation form and handler
	r.Get("/transactions/new", newTransactionForm(lineItemRepo))
	r.Post("/transactions", createTransaction(lineItemRepo))
}

func listTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Placeholder list – could be expanded.
		fmt.Fprintln(w, "Transactions endpoint – implementation pending.")
	}
}

// newTransactionForm renders a form to create a transaction.
func newTransactionForm(repo *db.InMemoryLineItemRepo) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("views/transactions/form.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		items := repo.GetAll()
		data := struct {
			Items []*domain.LineItem
		}{Items: items}
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// createTransaction processes the submitted transaction form.
func createTransaction(repo *db.InMemoryLineItemRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		lineItemID := r.FormValue("line-item-id")
		description := r.FormValue("description")
		amountStr := r.FormValue("amount")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		txn := &domain.Transaction{Description: description, Amount: amount}
		if err := repo.AddTransaction(lineItemID, txn); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Redirect back to transaction list or line item view.
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
	}
}
