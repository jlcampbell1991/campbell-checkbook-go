package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"html/template"
	"net/http"
)

// accountService provides HTTP handlers for account operations.
// It uses a template set for rendering HTML pages.
type accountService struct {
	repo db.AccountRepository
	tmpl *template.Template
}

// AccountRoutes returns a chi.Router with account‑related endpoints.
// Includes routes for listing, showing, creating, and deleting accounts.
func AccountRoutes(t *template.Template, repo db.AccountRepository) chi.Router {
	s := &accountService{repo: repo, tmpl: t}
	r := chi.NewRouter()
	// New account form and create endpoint.
	r.Get("/new", s.new)
	r.Post("/", s.create)
	// Existing routes (placeholder, can be expanded later).
	// r.Get("/{id}", s.show)
	// r.Delete("/{id}", s.delete)
	return r
}

// new renders the form for creating a new account.
func (s *accountService) new(w http.ResponseWriter, r *http.Request) {
	if err := s.tmpl.ExecuteTemplate(w, "accounts_add.tmpl", struct{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// create processes the submitted form and creates a new account.
func (s *accountService) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	acc := domain.Account{Name: name}
	if _, err := s.repo.Insert(r.Context(), acc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
