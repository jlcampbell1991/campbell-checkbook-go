package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook/internal/domain"
)

// RegisterAccounts registers account routes on the provided router.
func RegisterAccounts(r chi.Router, repo *db.InMemoryAccountRepo) {
	r.Get("/accounts", listAccounts(repo))
	r.Post("/accounts", createAccount(repo))
	r.Get("/accounts/{id}", showAccount(repo))
	r.Delete("/accounts/{id}", deleteAccount(repo))
}

func listAccounts(repo *db.InMemoryAccountRepo) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("views/accounts_list.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, _ := repo.List()
		tmpl.Execute(w, accounts)
	}
}

func createAccount(repo *db.InMemoryAccountRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "name required", http.StatusBadRequest)
			return
		}
		a := &domain.Account{Name: name}
		repo.Create(a)
		http.Redirect(w, r, fmt.Sprintf("/accounts/%d", a.ID), http.StatusSeeOther)
	}
}

func showAccount(repo *db.InMemoryAccountRepo) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("views/account_show.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		acc, err := repo.Get(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		tmpl.Execute(w, acc)
	}
}

func deleteAccount(repo *db.InMemoryAccountRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		err = repo.Delete(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	}
}
