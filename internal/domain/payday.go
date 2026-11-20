package domain

// PayDayId wraps a string identifier for a PayDay.
type PayDayId struct {
    Value string `json:"id"`
}

// PayDayReq represents the payload for creating or updating a PayDay.
type PayDayReq struct {
    Name      string  `json:"name"`
    Total     float64 `json:"total"`
    AccountId string  `json:"accountId"`
    // Deposits can be added later; omitted for brevity.
}

// PayDay is the persisted representation of a PayDay.
type PayDay struct {
    Id        PayDayId `json:"id"`
    Name      string   `json:"name"`
    Total     float64  `json:"total"`
    AccountId string   `json:"accountId"`
}
