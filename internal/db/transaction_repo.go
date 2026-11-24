package db

import (
	"context"
	"fmt"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// TransactionRepository defines the contract for transaction data access.
type TransactionRepository interface {
	// Index returns a page of transactions.
	Index(page int) ([]domain.Transaction, error)
	// GetAll returns all transactions for CSV download.
	GetAll() ([]domain.Transaction, error)
	// Create stores a new transaction and returns its generated ID.
	Create(ctx context.Context, tx domain.Transaction) (string, error)
	// Show returns a single transaction by ID.
	Show(ctx context.Context, id string) (domain.Transaction, error)
	// Update modifies an existing transaction.
	Update(ctx context.Context, tx domain.Transaction) error
	// Delete removes a transaction by ID.
	Delete(ctx context.Context, id string) error
}

// StubTransactionRepository provides a minimal in‑memory implementation used for compilation and testing.
type StubTransactionRepository struct{}

func (s *StubTransactionRepository) Index(page int) ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func (s *StubTransactionRepository) GetAll() ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func (s *StubTransactionRepository) Create(ctx context.Context, tx domain.Transaction) (string, error) {
	if tx.Memo == "" {
		return "", fmt.Errorf("memo is required")
	}
	// Return a deterministic placeholder ID.
	return "tx-123", nil
}

func (s *StubTransactionRepository) Show(ctx context.Context, id string) (domain.Transaction, error) {
	return domain.Transaction{ID: id}, nil
}

func (s *StubTransactionRepository) Update(ctx context.Context, tx domain.Transaction) error {
	return nil
}

func (s *StubTransactionRepository) Delete(ctx context.Context, id string) error {
	return nil
}
