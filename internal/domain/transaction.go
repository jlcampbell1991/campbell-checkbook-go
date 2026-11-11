package domain

// Transaction is a placeholder for transaction data associated with a line item.
type Transaction struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}
