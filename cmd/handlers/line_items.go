package handlers

import (
    "context"
    "fmt"
    "html/template"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

type lineItemService struct {
    repo db.LineItemRepository
    tmpl *template.Template
}

// LineItemRoutes returns a chi.Router with line‑item related endpoints.
func LineItemRoutes(t *template.Template, repo db.LineItemRepository) chi.Router {
    s := &lineItemService{repo: repo, tmpl: t}
    r := chi.NewRouter()
    r.Get("/{id}", s.show)
    return r
}

func (s *lineItemService) show(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")
    // Parse page query param, default to 0 on error or missing.
    page := 0
    if pStr := r.URL.Query().Get("page"); pStr != "" {
        if p, err := strconv.Atoi(pStr); err == nil && p >= 0 {
            page = p
        }
    }

    // Retrieve data via repository.
    lineItem, err := s.repo.Show(ctx, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    txs, err := s.repo.GetTransactions(ctx, id, page)
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
        LineItem    domain.LineItem
        Deps        []domain.PayDayDeposit
        Transactions []domain.Transaction
        NextURL     string
    }{
        LineItem:    lineItem,
        Deps:        deps,
        Transactions: txs,
        NextURL:     nextURL,
    }

    if err := s.tmpl.ExecuteTemplate(w, "line_items_show.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
