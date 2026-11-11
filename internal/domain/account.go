package domain

// Account represents a financial account.
// For now we only store an ID and a Name.
type Account struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
