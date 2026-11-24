package domain

// LineItem represents a financial line item belonging to an account.
type LineItem struct {
	ID        string  `json:"id"`
	AccountID string  `json:"accountId"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
}
