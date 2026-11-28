package domain

import (
    "fmt"
    "net/url"
    "strconv"
)

type PayDayId struct{ Value string }

type PayDayReq struct {
    Name      string  `json:"name"`
    Total     float64 `json:"total"`
    AccountId string  `json:"accountId"`
    // Deposits omitted for brevity
}

type PayDay struct {
    Id        PayDayId `json:"id"`
    Name      string   `json:"name"`
    Total     float64  `json:"total"`
    AccountId string   `json:"accountId"`
}

// PayDayFromForm parses URL form values into a PayDayReq.
func PayDayFromForm(form url.Values) (*PayDayReq, error) {
    name := form.Get("name")
    totalStr := form.Get("total")
    acct := form.Get("account")
    if name == "" || totalStr == "" || acct == "" {
        return nil, fmt.Errorf("missing required field")
    }
    total, err := strconv.ParseFloat(totalStr, 64)
    if err != nil {
        return nil, err
    }
    return &PayDayReq{Name: name, Total: total, AccountId: acct}, nil
}
