package handlers

import (
    "html/template"
    "net/http"
    "strconv"
    "github.com/go-chi/chi/v5"
)

// Transaction represents a simple financial transaction.
type Transaction struct {
    ID          int
    AccountID   string
    Description string
    Amount      float64
}

var (
    // In‑memory store for transactions. In a real app this would be a DB.
    transactions []Transaction
    nextTxnID    = 1
)

// RegisterTransactionRoutes registers the GET form and POST processing routes.
func RegisterTransactionRoutes(r chi.Router, tmpl *template.Template) {
    r.Get("/transactions/new", newTransactionFormHandler(tmpl))
    r.Post("/transactions", createTransactionHandler())
}

// newTransactionFormHandler renders the transaction creation form.
func newTransactionFormHandler(tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // No specific data required for now.
        if err := tmpl.ExecuteTemplate(w, "transactionForm", nil); err != nil {
            http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
        }
    }
}

// createTransactionHandler processes the submitted form and stores the transaction.
func createTransactionHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "invalid form data", http.StatusBadRequest)
            return
        }
        amountStr := r.FormValue("amount")
        amt, err := strconv.ParseFloat(amountStr, 64)
        if err != nil {
            http.Error(w, "invalid amount", http.StatusBadRequest)
            return
        }
        txn := Transaction{
            ID:          nextTxnID,
            AccountID:   r.FormValue("accountId"),
            Description: r.FormValue("description"),
            Amount:      amt,
        }
        nextTxnID++
        transactions = append(transactions, txn)
        // Redirect to home or a list page (not implemented).
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}
