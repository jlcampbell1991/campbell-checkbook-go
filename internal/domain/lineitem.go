package domain

type LineItem struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Balance   *float64 `json:"balance,omitempty"`
	Budget    *float64 `json:"budget,omitempty"`
	Active    bool     `json:"active"`
	AccountID string   `json:"accountId"`
}
