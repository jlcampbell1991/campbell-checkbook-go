package domain

type LineItem struct {
    ID        string  `json:"id"`
    AccountID string  `json:"accountId"`
    Name      string  `json:"name"`
    Balance   float64 `json:"balance"`
    // Additional fields can be added as needed.
}
