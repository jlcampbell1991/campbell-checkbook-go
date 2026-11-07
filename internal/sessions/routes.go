package sessions

import (
	"campbell-checkbook/internal/config"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers session/auth routes onto the provided router.
func RegisterRoutes(r chi.Router, cfg *config.Config) {
	r.Route("/sessions", func(r chi.Router) {
		r.Post("/login", loginHandler(cfg))
		r.Post("/logout", logoutHandler())
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(cfg))
			r.Get("/me", meHandler())
		})
	})
}
