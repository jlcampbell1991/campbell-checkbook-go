package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// RegisterTransactions registers transaction routes on the provided router.
func RegisterTransactions(r chi.Router) {
	// List transactions (placeholder implementation)
	r.Get("/transactions", listTransactions())
	// Additional transaction routes (create, edit, delete, etc.) can be added here.
}

func listTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Integrate with the transaction repository and proper template rendering.
		fmt.Fprintln(w, "Transactions endpoint – implementation pending.")
	}
}
