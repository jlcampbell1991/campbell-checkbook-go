package db

import (
    "context"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// AccountRepository defines the data access methods needed by the account handlers.
type AccountRepository interface {
    // Index returns all accounts.
    Index(ctx context.Context) ([]domain.Account, error)
    // Show returns a single account by ID.
    Show(ctx context.Context, id string) (domain.Account, error)
    // Insert creates a new account and returns its ID.
    Insert(ctx context.Context, a domain.Account) (string, error)
    // Update modifies an existing account.
    Update(ctx context.Context, a domain.Account, id string) error
    // Destroy removes an account.
    Destroy(ctx context.Context, id string) error
    // GetTransactions returns a page of transactions for the account.
    GetTransactions(ctx context.Context, id string, page int) ([]domain.Transaction, error)
    // GetDislocated returns transactions that are not linked to any account (or similar business rule).
    GetDislocated(ctx context.Context, id string) ([]domain.Transaction, error)
    // RelocateAll moves all dislocated transactions to the account.
    RelocateAll(ctx context.Context, id string) error
}

// StubAccountRepo is a minimal in‑memory implementation used to satisfy compilation.
// In a real application this would be backed by PostgreSQL via pgx.
type StubAccountRepo struct{}

func NewStubAccountRepo() AccountRepository { return &StubAccountRepo{} }

func (s *StubAccountRepo) Index(ctx context.Context) ([]domain.Account, error) { return nil, nil }
func (s *StubAccountRepo) Show(ctx context.Context, id string) (domain.Account, error) {
    // Return an empty account with the requested ID.
    return domain.Account{ID: id, Name: "Stub Account"}, nil
}
func (s *StubAccountRepo) Insert(ctx context.Context, a domain.Account) (string, error) { return a.ID, nil }
func (s *StubAccountRepo) Update(ctx context.Context, a domain.Account, id string) error { return nil }
func (s *StubAccountRepo) Destroy(ctx context.Context, id string) error { return nil }
func (s *StubAccountRepo) GetTransactions(ctx context.Context, id string, page int) ([]domain.Transaction, error) {
    return nil, nil
}
func (s *StubAccountRepo) GetDislocated(ctx context.Context, id string) ([]domain.Transaction, error) {
    return nil, nil
}
func (s *StubAccountRepo) RelocateAll(ctx context.Context, id string) error { return nil }
