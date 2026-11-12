package domain

// PayDayDeposit is a placeholder for payday deposit data linked to a line item.
type PayDayDeposit struct {
	ID         string  `json:"id"`
	PayDayID   string  `json:"pay_day_id"`
	LineItemID string  `json:"line_item_id"`
	Amount     float64 `json:"amount"`
}
