package db

import (
    "context"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// LineItemRepository defines data access methods for line items.
type LineItemRepository interface {
    // Show returns a single line item by its ID.
    Show(ctx context.Context, id string) (domain.LineItem, error)
    // GetTransactions returns a page of transactions for the given line item.
    GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error)
    // GetPayDayDeposits returns all pay‑day deposits linked to the line item.
    GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error)
}
