package lineitems

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers the line‑items CRUD endpoints onto the provided router.
func RegisterRoutes(r chi.Router, repo Repository) {
	r.Route("/line-items", func(r chi.Router) {
		r.Post("/", createLineItemHandler(repo))
		r.Get("/", listLineItemsHandler(repo))
		r.Get("/{id}", getLineItemHandler(repo))
		r.Put("/{id}", updateLineItemHandler(repo))
		r.Delete("/{id}", deleteLineItemHandler(repo))
	})
}
