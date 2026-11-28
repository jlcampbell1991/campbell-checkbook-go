package db

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// AccountRepository defines data access methods for accounts.
type AccountRepository interface {
	// Insert creates a new account and returns its ID.
	Insert(ctx context.Context, a domain.Account) (string, error)
	// Additional methods (Show, Delete, etc.) may be added later.
}

// StubAccountRepo provides a simple in‑memory stub implementation.
type StubAccountRepo struct{}

func (s *StubAccountRepo) Insert(ctx context.Context, a domain.Account) (string, error) {
	if a.ID == "" {
		a.ID = "stub-id"
	}
	return a.ID, nil
}

// Export a default repository instance that can be overridden by the application.
var AccountRepo AccountRepository = &StubAccountRepo{}
