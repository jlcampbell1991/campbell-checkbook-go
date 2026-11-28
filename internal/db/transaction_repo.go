package db

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// TransactionRepository defines data access methods for transactions.
type TransactionRepository interface {
	// Index returns a page of transactions. page is zero‑based.
	Index(page int) ([]domain.Transaction, error)
	// GetAll returns all transactions (e.g., for CSV download).
	GetAll() ([]domain.Transaction, error)
	// Create a new transaction, returning its ID.
	Create(ctx context.Context, tx domain.Transaction) (string, error)
	// Show retrieves a transaction by ID.
	Show(ctx context.Context, id string) (domain.Transaction, error)
	// Update modifies an existing transaction.
	Update(ctx context.Context, tx domain.Transaction) error
	// Delete removes a transaction by ID.
	Delete(ctx context.Context, id string) error
}

// Stub implementation returning empty data.
type transactionRepoStub struct{}

func (r *transactionRepoStub) Index(page int) ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func (r *transactionRepoStub) GetAll() ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func (r *transactionRepoStub) Create(ctx context.Context, tx domain.Transaction) (string, error) {
	return "stub-id", nil
}

func (r *transactionRepoStub) Show(ctx context.Context, id string) (domain.Transaction, error) {
	return domain.Transaction{}, nil
}

func (r *transactionRepoStub) Update(ctx context.Context, tx domain.Transaction) error { return nil }

func (r *transactionRepoStub) Delete(ctx context.Context, id string) error { return nil }

// Export a default repository instance (can be replaced later).
var TransactionRepo TransactionRepository = &transactionRepoStub{}

// NewStubTransactionRepo returns a stub implementation of TransactionRepository.
func NewStubTransactionRepo() TransactionRepository { return &transactionRepoStub{} }
