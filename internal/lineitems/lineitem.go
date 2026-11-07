package lineitems

import "time"

// LineItem represents a line item in a transaction.
type LineItem struct {
	ID          int64     `json:"id"`
	AccountID   int64     `json:"account_id"`
	Description string    `json:"description"`
	Amount      int64     `json:"amount"` // stored as smallest currency unit (e.g., cents)
	Currency    string    `json:"currency"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
