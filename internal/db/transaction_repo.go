package db

import "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"

type TransactionRepository interface {
	Index(page int) ([]domain.Transaction, error)
	GetAll() ([]domain.Transaction, error)
}

type StubTransactionRepo struct{}

func NewStubTransactionRepo() TransactionRepository { return &StubTransactionRepo{} }

func (s *StubTransactionRepo) Index(page int) ([]domain.Transaction, error) {
	// Stub implementation returns empty slice.
	return []domain.Transaction{}, nil
}

func (s *StubTransactionRepo) GetAll() ([]domain.Transaction, error) {
	// Stub implementation returns empty slice.
	return []domain.Transaction{}, nil
}
