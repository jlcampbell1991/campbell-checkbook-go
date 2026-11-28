package domain

// Account represents a financial account.
// Fields are kept simple for the initial implementation.
type Account struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Budget  float64 `json:"budget"`
}
