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

// RegisterTransfers registers transfer routes.
func RegisterTransfers(r chi.Router, lineItemRepo *db.InMemoryLineItemRepo) {
	r.Get("/transfers/new", newTransferForm(lineItemRepo))
	r.Post("/transfers", createTransfer(lineItemRepo))
}

// newTransferForm renders the transfer creation form.
func newTransferForm(repo *db.InMemoryLineItemRepo) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("views/transfers/form.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		// Get optional selected line-item-id query param
		selected := r.URL.Query().Get("line-item-id")
		// Retrieve all line items
		items := repo.GetAll()
		data := struct {
			SelectedID string
			Items      []*domain.LineItem
		}{
			SelectedID: selected,
			Items:      items,
		}
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// createTransfer processes the submitted transfer form.
func createTransfer(repo *db.InMemoryLineItemRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		memo := r.FormValue("memo")
		amountStr := r.FormValue("amount")
		toID := r.FormValue("to")
		fromID := r.FormValue("from")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		// Build Transfer struct (domain not yet defined, using simple struct)
		transfer := domain.Transfer{
			Memo:   memo,
			Amount: amount,
			To:     toID,
			From:   fromID,
		}
		// Placeholder: log the transfer; real implementation would create transactions.
		fmt.Fprintf(w, "Transfer received: %+v", transfer)
	}
}
