package db

import (
    "context"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// LineItemRepository defines data access methods for line items.
type LineItemRepository interface {
    Insert(ctx context.Context, li domain.LineItem) (string, error)
    // Show returns a single line item by its ID.
    Show(ctx context.Context, id string) (domain.LineItem, error)
    // GetTransactions returns a page of transactions for the given line item.
    GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error)
    // GetPayDayDeposits returns all pay‑day deposits linked to the line item.
    GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error)
}

// Stub implementation for development/testing.
type StubLineItemRepo struct{}

func (s *StubLineItemRepo) Show(ctx context.Context, id string) (domain.LineItem, error) {
    return domain.LineItem{}, nil
}

func (s *StubLineItemRepo) GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error) {
    return []domain.Transaction{}, nil
}

func (s *StubLineItemRepo) GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error) {
    return []domain.PayDayDeposit{}, nil
}

func (s *StubLineItemRepo) Insert(ctx context.Context, li domain.LineItem) (string, error) {
    if li.ID == "" {
        li.ID = "stub-id"
    }
    return li.ID, nil
}

// Export a default repository instance (can be replaced later).
var LineItemRepo LineItemRepository = &StubLineItemRepo{}
