package reconcile

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers the reconcile endpoints onto the provided router.
func RegisterRoutes(r chi.Router, repo Repository) {
	r.Route("/reconcile", func(r chi.Router) {
		r.Post("/jobs", createJobHandler(repo))
		r.Get("/jobs", listJobsHandler(repo))
		r.Get("/jobs/{id}", getJobHandler(repo))
	})
}
