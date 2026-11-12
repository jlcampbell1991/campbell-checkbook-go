package db

import (
	"errors"
	"fmt"
	"github.com/jlcampbell1991/campbell-checkbook/internal/domain"
	"strconv"
	"sync"
)

// InMemoryLineItemRepo provides a simple thread‑safe repository for line items.
type InMemoryLineItemRepo struct {
    // existing fields

	mu     sync.RWMutex
	items  map[string]*domain.LineItem
	nextID int64
	// For simplicity, transactions and deposits are stored in maps keyed by line item ID.
	txns     map[string][]*domain.Transaction
	deposits map[string][]*domain.PayDayDeposit
}

func NewInMemoryLineItemRepo() *InMemoryLineItemRepo {
	return &InMemoryLineItemRepo{
		items:    make(map[string]*domain.LineItem),
		txns:     make(map[string][]*domain.Transaction),
		deposits: make(map[string][]*domain.PayDayDeposit),
		nextID:   1,
	}
}

func (r *InMemoryLineItemRepo) Insert(item *domain.LineItem) (*domain.LineItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := strconv.FormatInt(r.nextID, 10)
	r.nextID++
	item.ID = id
	r.items[id] = item
	return item, nil
}

func (r *InMemoryLineItemRepo) Show(id string) (*domain.LineItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if it, ok := r.items[id]; ok {
		return it, nil
	}
	return nil, errors.New("line item not found")
}

func (r *InMemoryLineItemRepo) GetTransactions(id string, page int) ([]*domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Pagination ignored for stub implementation.
	if txs, ok := r.txns[id]; ok {
		return txs, nil
	}
	return []*domain.Transaction{}, nil
}

func (r *InMemoryLineItemRepo) GetAllByLineItem(id string) ([]*domain.PayDayDeposit, error) {
    // existing method unchanged

    // existing method unchanged

	r.mu.RLock()
	defer r.mu.RUnlock()
	if deps, ok := r.deposits[id]; ok {
		return deps, nil
	}
	return []*domain.PayDayDeposit{}, nil
}

// GetAll returns all line items.
func (r *InMemoryLineItemRepo) GetAll() []*domain.LineItem {
    r.mu.RLock()
    defer r.mu.RUnlock()
    result := make([]*domain.LineItem, 0, len(r.items))
    for _, v := range r.items {
        result = append(result, v)
    }
    return result
}

// AddTransaction adds a transaction to the specified line item.
func (r *InMemoryLineItemRepo) AddTransaction(lineItemID string, txn *domain.Transaction) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    // Ensure the line item exists.
    if _, ok := r.items[lineItemID]; !ok {
        return fmt.Errorf("line item %s not found", lineItemID)
    }
    // Simple ID generation.
    txn.ID = strconv.Itoa(len(r.txns[lineItemID]) + 1)
    r.txns[lineItemID] = append(r.txns[lineItemID], txn)
    return nil
}

