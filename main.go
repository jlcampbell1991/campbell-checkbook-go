package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jlcampbell1991/campbell-checkbook/cmd/handlers"
	"github.com/jlcampbell1991/campbell-checkbook/internal/db"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	repo := db.NewInMemoryAccountRepo()

	// Register account routes
	handlers.RegisterAccounts(r, repo)

	// Register line‑item routes (using its own in‑memory repository)
	handlers.RegisterTransactions(r)
	lineItemRepo := db.NewInMemoryLineItemRepo()
	handlers.RegisterLineItems(r, lineItemRepo)

	// Serve static files (if any)
	// r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("listen error: %v", err)
	}
}
