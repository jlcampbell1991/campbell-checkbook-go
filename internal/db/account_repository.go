package db

import (
	"context"
	"fmt"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// AccountRepository defines the contract for persisting accounts.
type AccountRepository interface {
	// Insert stores a new account and returns its generated ID.
	Insert(ctx context.Context, a domain.Account) (string, error)
}

// StubAccountRepository provides a minimal in‑memory implementation used for compilation and testing.
type StubAccountRepository struct{}

func (s *StubAccountRepository) Insert(ctx context.Context, a domain.Account) (string, error) {
	// In a real implementation, this would generate a UUID or use DB auto‑generated ID.
	// For the stub we simply return a placeholder ID.
	if a.Name == "" {
		return "", fmt.Errorf("name is required")
	}
	// Return a deterministic fake ID.
	return "acc-123", nil
}
