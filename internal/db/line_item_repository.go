package db

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// LineItemRepository defines the contract for line‑item data access.
type LineItemRepository interface {
	// Show returns a single line‑item by its ID.
	Show(ctx context.Context, id string) (domain.LineItem, error)
	// GetTransactions returns a page of transactions for the line‑item.
	GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error)
	// GetPayDayDeposits returns all payday deposits linked to the line‑item.
	GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error)
}

// StubLineItemRepository provides a minimal in‑memory implementation used for compilation and testing.
type StubLineItemRepository struct{}

func (s *StubLineItemRepository) Show(ctx context.Context, id string) (domain.LineItem, error) {
	// Return a placeholder line‑item; real implementation will query the database.
	return domain.LineItem{ID: id, AccountID: "acc-1", Name: "Sample Line Item", Balance: 0.0}, nil
}

func (s *StubLineItemRepository) GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error) {
	// Return an empty slice for now.
	return []domain.Transaction{}, nil
}

func (s *StubLineItemRepository) GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error) {
	// Return an empty slice for now.
	return []domain.PayDayDeposit{}, nil
}
