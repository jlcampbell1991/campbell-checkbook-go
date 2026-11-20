package domain

import "time"

type Transaction struct {
    ID        string    `json:"id"`
    AccountID string    `json:"account_id"`
    Amount    float64   `json:"amount"`
    Category  string    `json:"category"` // e.g., "Deposit" or "Withdrawal"
    Date      time.Time `json:"date"`
    // Additional fields can be added as needed.
}
