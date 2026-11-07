package main

import (
	"fmt"
	"log"
	"net/http"

	"campbell-checkbook/internal/accounts"
	"campbell-checkbook/internal/config"
	"campbell-checkbook/internal/db"
	"campbell-checkbook/internal/lineitems"
	"campbell-checkbook/internal/paydays"
	"campbell-checkbook/internal/reconcile"
	"campbell-checkbook/internal/server"
	"campbell-checkbook/internal/sessions"
	"campbell-checkbook/internal/transactions"

	"github.com/subosito/gotenv"
)

func main() {
	// Load .env if present (ignore error if file missing)
	_ = gotenv.Load(".env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize DB (optional, for now just test connection)
	dbConn, err := db.New(cfg)
	if err != nil {
		log.Printf("warning: could not connect to database: %v", err)
		// Continue without DB for now; health endpoint still works.
	}

	router := server.NewRouter()
	// Register domain routes
	if dbConn != nil {
		// Accounts
		accountsRepo := accounts.NewPGRepository(dbConn)
		accounts.RegisterRoutes(router, accountsRepo)
		// Line items
		lineItemsRepo := lineitems.NewPGRepository(dbConn)
		lineitems.RegisterRoutes(router, lineItemsRepo)
		// Transactions
		transactionsRepo := transactions.NewPGRepository(dbConn)
		transactions.RegisterRoutes(router, transactionsRepo)
		// Paydays
		paydaysRepo := paydays.NewPGRepository(dbConn)
		paydays.RegisterRoutes(router, paydaysRepo)
		// Reconcile
		reconcileRepo := reconcile.NewPGRepository(dbConn)
		reconcile.RegisterRoutes(router, reconcileRepo)
	}
	// Register session routes
	sessions.RegisterRoutes(router, cfg)

	addr := fmt.Sprintf(":" + cfg.Port)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
