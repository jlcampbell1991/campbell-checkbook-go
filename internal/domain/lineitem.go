package domain

// LineItem represents a budget line item.
type LineItem struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Budget    *float64 `json:"budget,omitempty"`
	Balance   float64  `json:"balance"`
	Active    bool     `json:"active"`
	AccountID string   `json:"account_id"`
}
