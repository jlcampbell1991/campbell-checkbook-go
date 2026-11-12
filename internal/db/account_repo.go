package db

import (
    "errors"
    "sync"
    "campbell-checkbook/internal/domain"
)

// AccountRepository defines storage operations for accounts.
type AccountRepository interface {
    Index() ([]domain.Account, error)
    Insert(name string) (domain.Account, error)
    Show(id string) (domain.Account, error)
    Update(id string, name string) (domain.Account, error)
    Delete(id string) error
}

// memAccountRepo is a simple in‑memory implementation used for the form demo.
type memAccountRepo struct {
    mu       sync.RWMutex
    accounts map[string]domain.Account
}

func NewMemAccountRepo() AccountRepository {
    return &memAccountRepo{accounts: make(map[string]domain.Account)}
}

func (r *memAccountRepo) Index() ([]domain.Account, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    list := make([]domain.Account, 0, len(r.accounts))
    for _, a := range r.accounts {
        list = append(list, a)
    }
    return list, nil
}

func (r *memAccountRepo) Insert(name string) (domain.Account, error) {
    if name == "" {
        return domain.Account{}, errors.New("name cannot be empty")
    }
    a := domain.NewAccount(name)
    r.mu.Lock()
    r.accounts[a.ID] = a
    r.mu.Unlock()
    return a, nil
}

func (r *memAccountRepo) Show(id string) (domain.Account, error) {
    r.mu.RLock()
    a, ok := r.accounts[id]
    r.mu.RUnlock()
    if !ok {
        return domain.Account{}, errors.New("account not found")
    }
    return a, nil
}

func (r *memAccountRepo) Update(id string, name string) (domain.Account, error) {
    if name == "" {
        return domain.Account{}, errors.New("name cannot be empty")
    }
    r.mu.Lock()
    defer r.mu.Unlock()
    a, ok := r.accounts[id]
    if !ok {
        return domain.Account{}, errors.New("account not found")
    }
    a.Name = name
    r.accounts[id] = a
    return a, nil
}

func (r *memAccountRepo) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.accounts[id]; !ok {
        return errors.New("account not found")
    }
    delete(r.accounts, id)
    return nil
}
