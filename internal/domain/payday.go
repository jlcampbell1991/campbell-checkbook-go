package domain

// PayDayId represents the unique identifier for a PayDay.
type PayDayId struct {
    Value string `json:"value"`
}

// PayDayReq represents the request payload for creating or updating a PayDay.
type PayDayReq struct {
    Name      string  `json:"name"`
    Total     float64 `json:"total"`
    AccountId string  `json:"accountId"`
    // Deposits omitted for brevity; can be added later.
}

// PayDay represents a pay day entity.
type PayDay struct {
    Id        PayDayId `json:"id"`
    Name      string   `json:"name"`
    Total     float64  `json:"total"`
    AccountId string   `json:"accountId"`
}
