package transactions

import "time"

// Transaction represents a financial transaction linking an account and a line item.
type Transaction struct {
	ID          int64     `json:"id"`
	AccountID   int64     `json:"account_id"`
	LineItemID  int64     `json:"line_item_id"`
	Amount      int64     `json:"amount"` // stored as smallest currency unit (e.g., cents)
	Currency    string    `json:"currency"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
