package domain

// Transfer represents a money transfer between two line items.
type Transfer struct {
    Memo   string  `json:"memo"`
    Amount float64 `json:"amount"`
    To     string  `json:"to"`
    From   string  `json:"from"`
}
