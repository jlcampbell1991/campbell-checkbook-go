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
    // Create inserts a new transaction and returns its ID.
    Create(ctx context.Context, tx domain.Transaction) (string, error)
    // Show retrieves a transaction by ID.
    Show(ctx context.Context, id string) (domain.Transaction, error)
    // Update modifies an existing transaction.
    Update(ctx context.Context, tx domain.Transaction) error
    // Delete removes a transaction by ID.
    Delete(ctx context.Context, id string) error
}

// Stub implementation returning empty data.
type StubTransactionRepo struct{}

func (r *StubTransactionRepo) Index(page int) ([]domain.Transaction, error) {
    return []domain.Transaction{}, nil
}

func (r *StubTransactionRepo) GetAll() ([]domain.Transaction, error) {
    return []domain.Transaction{}, nil
}

func (r *StubTransactionRepo) Create(ctx context.Context, tx domain.Transaction) (string, error) {
    // In stub, assign a dummy ID if not set.
    if tx.ID == "" {
        tx.ID = "stub-id"
    }
    return tx.ID, nil
}

func (r *StubTransactionRepo) Show(ctx context.Context, id string) (domain.Transaction, error) {
    return domain.Transaction{ID: id}, nil
}

func (r *StubTransactionRepo) Update(ctx context.Context, tx domain.Transaction) error {
    return nil
}

func (r *StubTransactionRepo) Delete(ctx context.Context, id string) error {
    return nil
}

// Export a default repository instance (can be replaced later).
var TransactionRepo TransactionRepository = &StubTransactionRepo{}
