package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Account struct {
	ID      int
	Name    string
	Balance float64
}

// RegisterAccountRoutes registers routes for account show and delete.
func RegisterAccountRoutes(r chi.Router, tmpl *template.Template) {
	r.Get("/accounts/{id}", showAccountHandler(tmpl))
	r.Post("/accounts/{id}/delete", deleteAccountHandler())
}

// showAccountHandler returns a handler that renders the account show page using the "accountShow" template.
func showAccountHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		// Dummy data – in real implementation fetch from DB
		data := struct {
			Account struct {
				ID   string
				Name string
			}
			LineItems []struct {
				Name    string
				Balance float64
			}
		}{}
		data.Account.ID = id
		data.Account.Name = fmt.Sprintf("Account %s", id)
		// Example line item
		data.LineItems = []struct {
			Name    string
			Balance float64
		}{{Name: "Sample Item", Balance: 10.0}}

		// Execute the accountShow template into the layout
		// First render accountShow into a buffer
		var buf strings.Builder
		if err := tmpl.ExecuteTemplate(&buf, "accountShow", data); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Wrap with application layout
		pageData := struct {
			Title string
			Body  template.HTML
		}{
			Title: fmt.Sprintf("Account %s", id),
			Body:  template.HTML(buf.String()),
		}
		if err := tmpl.ExecuteTemplate(w, "application.tmpl", pageData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// deleteAccountHandler returns a handler that processes account deletion.
func deleteAccountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real implementation, delete from DB.
		// For now, just redirect to home.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
