package handlers

import (
        "fmt"
    "html/template"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

type accountService struct {
    repo db.AccountRepository
    tmpl *template.Template
}

// AccountRoutes returns a chi.Router with account related endpoints.
func AccountRoutes(t *template.Template, repo db.AccountRepository) chi.Router {
    h := &accountService{repo: repo, tmpl: t}
    r := chi.NewRouter()
    r.Get("/{id}", h.show)
    r.Get("/{id}/delete", h.delete)
    return r
}

func (h *accountService) show(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")
    // Parse page query param, default to 0
    page := 0
    if pStr := r.URL.Query().Get("page"); pStr != "" {
        if p, err := strconv.Atoi(pStr); err == nil && p >= 0 {
            page = p
        }
    }

    account, err := h.repo.Show(ctx, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    transactions, _ := h.repo.GetTransactions(ctx, id, page)
    dislocated, _ := h.repo.GetDislocated(ctx, id)

    var dlTotal float64
    for _, dl := range dislocated {
        if dl.Category == "Withdrawal" {
            dlTotal -= dl.Amount
        } else {
            dlTotal += dl.Amount
        }
    }
    nextURL := fmt.Sprintf("/accounts/%s?page=%d", id, page+1)
    data := struct {
        Account      domain.Account
        Transactions []domain.Transaction
        Dislocated   []domain.Transaction
        DlTotal      float64
        NextURL      string
    }{
        Account:      account,
        Transactions: transactions,
        Dislocated:   dislocated,
        DlTotal:      dlTotal,
        NextURL:      nextURL,
    }
    if err := h.tmpl.ExecuteTemplate(w, "accounts_show.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (h *accountService) delete(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")
    if err := h.repo.Destroy(ctx, id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
