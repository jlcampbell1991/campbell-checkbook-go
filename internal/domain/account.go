package domain

type Account struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Balance   *float64   `json:"balance,omitempty"`
	Budget    *float64   `json:"budget,omitempty"`
	LineItems []LineItem `json:"lineItems,omitempty"`
}
