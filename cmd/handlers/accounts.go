package handlers

import (
    "html/template"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

// RegisterAccountsRoutes registers the account related routes.
func RegisterAccountsRoutes(r chi.Router, tmpl *template.Template, repo *db.InMemoryAccountRepo) {
    // Show form for new account
    r.Get("/accounts/new", func(w http.ResponseWriter, r *http.Request) {
        data := struct {
            Title string
        }{Title: "New Account"}
        if err := tmpl.ExecuteTemplate(w, "form.tmpl", data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })

    // Handle form submission
    r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "invalid form", http.StatusBadRequest)
            return
        }
        name := r.PostFormValue("name")
        if name == "" {
            http.Error(w, "name required", http.StatusBadRequest)
            return
        }
        repo.Insert(name)
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })
}
