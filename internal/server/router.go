package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// NewRouter returns a chi.Router with common middleware and routes.
func NewRouter() chi.Router {
	r := chi.NewRouter()
	// Basic logger
	r.Use(middleware.Logger)
	// Recoverer to handle panics
	r.Use(middleware.Recoverer)

	// Health endpoint
	r.Get("/health", healthHandler)
	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
