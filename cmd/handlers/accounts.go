package handlers

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// AccountRepository abstracts the data access for accounts.
type AccountRepository interface {
	Index(ctx context.Context) ([]domain.Account, error)
}
type AccountsHandler struct {
	repo AccountRepository
	tmpl *template.Template
}

// NewAccountsHandler creates a new handler with the given repository and parsed templates.
func NewAccountsHandler(repo AccountRepository, tmpl *template.Template) *AccountsHandler {
	return &AccountsHandler{repo: repo, tmpl: tmpl}
}

// Router returns a chi router with the accounts routes.
func (h *AccountsHandler) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.index)
	return r
}

// index handles GET "/" and renders the accounts index page.
func (h *AccountsHandler) index(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.repo.Index(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	data := struct {
		Title    string
		Accounts []domain.Account
		Body     template.HTML
	}{
		Title:    "Home",
		Accounts: accounts,
	}
	// Render the body template into a buffer.
	var bodyBuf bytes.Buffer
	if err := h.tmpl.ExecuteTemplate(&bodyBuf, "body", data); err != nil {
		log.Printf("body template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Body = template.HTML(bodyBuf.String())
	// Render using the application layout template.
	if err := h.tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
