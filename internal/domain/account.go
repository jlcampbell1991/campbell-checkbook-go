package domain

// Account represents a financial account.
// Fields are kept minimal for the demo implementation.
type Account struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
