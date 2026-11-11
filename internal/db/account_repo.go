package db

import (
	"errors"
	"github.com/jlcampbell1991/campbell-checkbook/internal/domain"
	"sync"
)

// InMemoryAccountRepo is a simple thread‑safe repository used for the demo.
type InMemoryAccountRepo struct {
	mu       sync.RWMutex
	accounts map[int64]*domain.Account
	nextID   int64
}

func NewInMemoryAccountRepo() *InMemoryAccountRepo {
	return &InMemoryAccountRepo{
		accounts: make(map[int64]*domain.Account),
		nextID:   1,
	}
}

func (r *InMemoryAccountRepo) Create(a *domain.Account) (*domain.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	a.ID = r.nextID
	r.nextID++
	r.accounts[a.ID] = a
	return a, nil
}

func (r *InMemoryAccountRepo) Get(id int64) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if a, ok := r.accounts[id]; ok {
		return a, nil
	}
	return nil, errors.New("account not found")
}

func (r *InMemoryAccountRepo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.accounts[id]; ok {
		delete(r.accounts, id)
		return nil
	}
	return errors.New("account not found")
}

func (r *InMemoryAccountRepo) List() ([]*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*domain.Account, 0, len(r.accounts))
	for _, a := range r.accounts {
		list = append(list, a)
	}
	return list, nil
}
