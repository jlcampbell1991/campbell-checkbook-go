package db

import "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"

// TransactionRepository defines data access methods for transactions.
type TransactionRepository interface {
	// Index returns a page of transactions. page is zero‑based.
	Index(page int) ([]domain.Transaction, error)
	// GetAll returns all transactions (e.g., for CSV download).
	GetAll() ([]domain.Transaction, error)
}

// Stub implementation returning empty data.
type transactionRepoStub struct{}

func (r *transactionRepoStub) Index(page int) ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func (r *transactionRepoStub) GetAll() ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

// Export a default repository instance (can be replaced later).
var TransactionRepo TransactionRepository = &transactionRepoStub{}
