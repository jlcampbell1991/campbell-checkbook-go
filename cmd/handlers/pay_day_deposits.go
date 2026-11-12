package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook/internal/domain"
	"html/template"
	"log"
	"net/http"
)

// RegisterPayDayDeposits registers routes for pay‑day deposit related pages.
func RegisterPayDayDeposits(r chi.Router, lineItemRepo *db.InMemoryLineItemRepo) {
	// New deposit form
	r.Get("/pay-day-deposits/new", newPayDayDepositForm(lineItemRepo))
	// Additional routes (edit, create, delete) could be added here.
}

// newPayDayDepositForm renders the form for creating a new PayDayDeposit.
func newPayDayDepositForm(lineItemRepo *db.InMemoryLineItemRepo) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("views/pay_day_deposits/form.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		// Query parameters for pre‑selection (optional).
		maybeItemId := r.URL.Query().Get("itemId")
		maybePayDayId := r.URL.Query().Get("payDayId")

		// Gather data needed for the form.
		lineItems := lineItemRepo.GetAll()
		// For now, PayDays are empty as the PayDay repository is not yet implemented.
		payDays := []domain.PayDay{}

		data := struct {
			Title         string
			Action        string
			Method        string
			MaybeDep      *domain.PayDayDeposit
			PayDays       []domain.PayDay
			Items         []*domain.LineItem
			MaybeItemId   string
			MaybePayDayId string
		}{
			Title:         "New Pay Day Deposit",
			Action:        "/pay-day-deposits",
			Method:        "post",
			MaybeDep:      nil,
			PayDays:       payDays,
			Items:         lineItems,
			MaybeItemId:   maybeItemId,
			MaybePayDayId: maybePayDayId,
		}
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
