package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// accountService provides HTTP handlers for account creation.
type accountService struct {
	repo db.AccountRepository
	tmpl *template.Template
}

// AccountRoutes registers the account‑creation routes.
// It mounts under the "/accounts" path by the caller.
func AccountRoutes(t *template.Template, repo db.AccountRepository) chi.Router {
	h := &accountService{repo: repo, tmpl: t}
	r := chi.NewRouter()
	// Render the new‑account form.
	r.Get("/new", h.new)
	// Handle form submission.
	r.Post("/", h.create)
	return r
}

// new renders the account creation form.
func (h *accountService) new(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "accounts_add.tmpl", struct {
		Title string
		Body  string
	}{Title: "", Body: ""}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// create processes the posted form and creates a new account.
func (h *accountService) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	acc := domain.Account{Name: name}
	if _, err := h.repo.Insert(r.Context(), acc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
