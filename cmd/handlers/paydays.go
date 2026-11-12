package handlers

import (
    "html/template"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/jlcampbell1991/campbell-checkbook/internal/db"
    "github.com/jlcampbell1991/campbell-checkbook/internal/domain"
)

// RegisterPayDays registers routes for pay‑day related pages.
func RegisterPayDays(r chi.Router, lineItemRepo *db.InMemoryLineItemRepo, accountRepo *db.InMemoryAccountRepo) {
    // Show form for creating a new pay day.
    r.Get("/pay-days/new", newPayDayForm(lineItemRepo, accountRepo))
    // Future routes (create, edit, etc.) would be added here.
}

// newPayDayForm renders the pay‑day creation form.
func newPayDayForm(lineItemRepo *db.InMemoryLineItemRepo, accountRepo *db.InMemoryAccountRepo) http.HandlerFunc {
    tmpl := template.Must(template.ParseFiles("views/pay_days/form.html"))
    return func(w http.ResponseWriter, r *http.Request) {
        // Gather data needed for the form.
        lineItems := lineItemRepo.GetAll()
        accounts, err := accountRepo.List()
        if err != nil {
            log.Printf("error listing accounts: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        // Prepare a map of deposits (empty for a new form).
        deposits := make(map[string]float64)
        // Empty PayDay struct for default values.
        data := struct {
            Title     string
            Action    string
            Method    string
            PayDay    struct {
                Name      string
                Total     float64
                AccountID int64
            }
            Accounts  []*domain.Account
            LineItems []*domain.LineItem
            Deposits  map[string]float64
        }{
            Title:    "New Pay Day",
            Action:   "/pay-days",
            Method:   "post",
            Accounts: accounts,
            LineItems: lineItems,
            Deposits: deposits,
        }
        if err := tmpl.Execute(w, data); err != nil {
            log.Printf("template error: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    }
}
