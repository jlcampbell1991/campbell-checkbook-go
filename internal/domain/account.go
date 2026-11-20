package domain

type Account struct {
    ID      string   `json:"id"`
    Name    string   `json:"name"`
    Balance float64 `json:"balance"`
    Budget  float64 `json:"budget"`
    // Additional fields can be added as needed.
}
