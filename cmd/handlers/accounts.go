package handlers

import (
    "net/http"
    "html/template"
    "github.com/go-chi/chi/v5"
    "campbell-checkbook/internal/db"
    "campbell-checkbook/internal/domain"
)

// RegisterAccounts registers account related routes on the provided router.
func RegisterAccounts(r chi.Router, repo db.AccountRepository) {
    r.Get("/accounts/new", func(w http.ResponseWriter, r *http.Request) {
        renderForm(w, "New Account", "/accounts", "post", nil)
    })

    r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "invalid form", http.StatusBadRequest)
            return
        }
        name := r.FormValue("name")
        if _, err := repo.Insert(name); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })
}

// renderForm renders the account form template.
func renderForm(w http.ResponseWriter, title, action, method string, maybeAcc *domain.Account) {
    tmpl := template.Must(template.ParseFiles("../../views/accounts/form.tmpl"))
    data := struct {
        Title   string
        Action  string
        Method  string
        Account *domain.Account
    }{Title: title, Action: action, Method: method, Account: maybeAcc}
    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
