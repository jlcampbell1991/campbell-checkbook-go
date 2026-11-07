package paydays

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers payday and deposit endpoints onto the provided router.
func RegisterRoutes(r chi.Router, repo Repository) {
	// PayDay CRUD
	r.Route("/paydays", func(r chi.Router) {
		r.Post("/", createPayDayHandler(repo))
		r.Get("/", listPayDaysHandler(repo))
		r.Get("/{id}", getPayDayHandler(repo))
		r.Put("/{id}", updatePayDayHandler(repo))
		r.Delete("/{id}", deletePayDayHandler(repo))
		// Deposits nested under a specific payday
		r.Route("/{id}/deposits", func(r chi.Router) {
			r.Post("/", createDepositHandler(repo))
			r.Get("/", listDepositsHandler(repo))
			r.Get("/{deposit_id}", getDepositHandler(repo))
			r.Put("/{deposit_id}", updateDepositHandler(repo))
			r.Delete("/{deposit_id}", deleteDepositHandler(repo))
		})
	})
}
