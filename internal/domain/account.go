package domain

import "github.com/google/uuid"

type Account struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// NewAccount creates a new Account with a generated UUID.
func NewAccount(name string) Account {
    return Account{ID: uuid.NewString(), Name: name}
}
