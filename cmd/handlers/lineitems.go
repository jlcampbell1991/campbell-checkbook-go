package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook/internal/domain"
	"html/template"
	"net/http"
	"strconv"
)

// RegisterLineItems registers line‑item routes.
func RegisterLineItems(r chi.Router, repo *db.InMemoryLineItemRepo) {
	r.Get("/line-items/{id}", showLineItem(repo))
	// Additional routes could be added here.
}

func showLineItem(repo *db.InMemoryLineItemRepo) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("views/line_items_show.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		pageStr := r.URL.Query().Get("page")
		page := 0
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil {
				page = p
			}
		}
		// Fetch data (stub implementations).
		lineItem, err := repo.Show(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		transactions, _ := repo.GetTransactions(id, page)
		deps, _ := repo.GetAllByLineItem(id)
		nextURL := "/line-items/" + id + "?page=" + strconv.Itoa(page+1)
		data := struct {
			LineItem     *domain.LineItem
			Deps         []*domain.PayDayDeposit
			Transactions []*domain.Transaction
			NextURL      string
		}{
			LineItem:     lineItem,
			Deps:         deps,
			Transactions: transactions,
			NextURL:      nextURL,
		}
		tmpl.Execute(w, data)
	}
}
