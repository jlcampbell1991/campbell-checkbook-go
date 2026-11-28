package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

type lineItemService struct {
	repo db.LineItemRepository
	tmpl *template.Template
}

// LineItemRoutes returns a chi.Router handling line‑item related endpoints.
func LineItemRoutes(t *template.Template, repo db.LineItemRepository) chi.Router {
	s := &lineItemService{repo: repo, tmpl: t}
	r := chi.NewRouter()
	r.Get("/{id}", s.show)
	return r
}

func (s *lineItemService) show(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	// Parse page query parameter, default to 0
	page := 0
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v >= 0 {
			page = v
		}
	}

	// Retrieve data from repository
	lineItem, err := s.repo.Show(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	transactions, err := s.repo.GetTransactions(ctx, id, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	deps, err := s.repo.GetPayDayDeposits(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nextURL := fmt.Sprintf("/line-items/%s?page=%d", id, page+1)

	data := struct {
		LineItem     domain.LineItem
		Deps         []domain.PayDayDeposit
		Transactions []domain.Transaction
		NextURL      string
	}{
		LineItem:     lineItem,
		Deps:         deps,
		Transactions: transactions,
		NextURL:      nextURL,
	}

	if err := s.tmpl.ExecuteTemplate(w, "line_items_show.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
