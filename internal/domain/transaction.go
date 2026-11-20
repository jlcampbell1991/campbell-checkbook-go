package domain

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	Date            time.Time `json:"date"`
	Memo            string    `json:"memo"`
	Amount          float64   `json:"amount"`
	Category        string    `json:"category"`
	LineItemID      string    `json:"line_item_id"`
	LineItem        string    `json:"line_item"`
	LineItemBalance float64   `json:"line_item_balance"`
	Reconciled      bool      `json:"reconciled"`
	AccountID       string    `json:"account_id"`
	LocationID      string    `json:"location_id"`
}
