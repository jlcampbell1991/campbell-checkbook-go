package db

import (
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// Index returns a page of transactions. Stub implementation returns empty slice.
func Index(page int) ([]domain.Transaction, error) {
	// TODO: replace with real DB query using pagination (pageSize = 20)
	return []domain.Transaction{}, nil
}

// GetAll returns all transactions for CSV download. Stub implementation.
func GetAll() ([]domain.Transaction, error) {
	// TODO: replace with real DB query to fetch all transactions.
	return []domain.Transaction{}, nil
}
