package domain

// PayDayDeposit represents a deposit made on a payday linked to a line item.
type PayDayDeposit struct {
    ID         string  `json:"id"`
    LineItemID string  `json:"lineItemId"`
    Amount     float64 `json:"amount"`
    Date       string  `json:"date"`
}
