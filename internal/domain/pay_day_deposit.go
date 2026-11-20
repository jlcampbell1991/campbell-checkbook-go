package domain

type PayDayDeposit struct {
	ID         string  `json:"id"`
	LineItemID string  `json:"lineItemId"`
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	// Additional fields can be added as needed.
}
