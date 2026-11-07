package transactions

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers transaction endpoints onto the provided router.
func RegisterRoutes(r chi.Router, repo Repository) {
	r.Route("/transactions", func(r chi.Router) {
		// Existing CRUD routes could be added later.
		r.Post("/import", importCSVHandler(repo))
	})
}
