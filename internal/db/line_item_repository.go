package db

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// LineItemRepository defines data access methods needed by the line‑item handlers.
type LineItemRepository interface {
	Show(ctx context.Context, id string) (domain.LineItem, error)
	GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error)
	GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error)
}

// StubLineItemRepo is a minimal in‑memory implementation used to satisfy compilation.
// In a real application this would be backed by PostgreSQL via pgx.
type StubLineItemRepo struct{}

func NewStubLineItemRepo() LineItemRepository { return &StubLineItemRepo{} }

func (s *StubLineItemRepo) Show(ctx context.Context, id string) (domain.LineItem, error) {
	// Return a placeholder line item.
	return domain.LineItem{ID: id, AccountID: "acct-1", Name: "Stub Line Item", Balance: 0.0}, nil
}

func (s *StubLineItemRepo) GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error) {
	// Return empty slice for now.
	return nil, nil
}

func (s *StubLineItemRepo) GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error) {
	// Return empty slice for now.
	return nil, nil
}
