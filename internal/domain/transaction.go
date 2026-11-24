package domain

// Transaction represents a financial transaction associated with a line item.
type Transaction struct {
    ID         string  `json:"id"`
    LineItemID string  `json:"lineItemId"`
    Amount     float64 `json:"amount"`
    Category   string  `json:"category"`
    Date       string  `json:"date"`
}
