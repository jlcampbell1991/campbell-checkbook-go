package paydays

import "time"

// PayDay represents a scheduled payday.
type PayDay struct {
	ID        int64     `json:"id"`
	Date      time.Time `json:"date"`
	Amount    int64     `json:"amount"` // smallest currency unit
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PayDayDeposit represents a deposit of a payday into an account.
type PayDayDeposit struct {
	ID        int64     `json:"id"`
	PayDayID  int64     `json:"payday_id"`
	AccountID int64     `json:"account_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
