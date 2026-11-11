package db

import (
    "fmt"
    "sync"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// InMemoryAccountRepo provides a simple thread‑safe repository for accounts.
type InMemoryAccountRepo struct {
    mu       sync.Mutex
    accounts map[string]domain.Account
    nextID   int
}

// NewInMemoryAccountRepo creates a new repository instance.
func NewInMemoryAccountRepo() *InMemoryAccountRepo {
    return &InMemoryAccountRepo{accounts: make(map[string]domain.Account)}
}

// Insert adds a new account and returns the created account.
func (r *InMemoryAccountRepo) Insert(name string) domain.Account {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.nextID++
    id := fmt.Sprintf("%d", r.nextID)
    acc := domain.Account{ID: id, Name: name}
    r.accounts[id] = acc
    return acc
}

// List returns all accounts.
func (r *InMemoryAccountRepo) List() []domain.Account {
    r.mu.Lock()
    defer r.mu.Unlock()
    result := make([]domain.Account, 0, len(r.accounts))
    for _, a := range r.accounts {
        result = append(result, a)
    }
    return result
}
