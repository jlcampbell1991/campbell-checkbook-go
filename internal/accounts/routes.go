package accounts

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers the accounts CRUD endpoints onto the provided router.
func RegisterRoutes(r chi.Router, repo Repository) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", createAccountHandler(repo))
		r.Get("/", listAccountsHandler(repo))
		r.Get("/{id}", getAccountHandler(repo))
		r.Put("/{id}", updateAccountHandler(repo))
		r.Delete("/{id}", deleteAccountHandler(repo))
	})
}
