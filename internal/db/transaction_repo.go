package db

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

type TransactionRepository interface {
	Index(page int) ([]domain.Transaction, error)
	GetAll() ([]domain.Transaction, error)
	Create(ctx context.Context, tx domain.Transaction) (string, error)
	Show(ctx context.Context, id string) (domain.Transaction, error)
	Update(ctx context.Context, tx domain.Transaction) error
	Delete(ctx context.Context, id string) error
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

func (s *StubTransactionRepo) Create(ctx context.Context, tx domain.Transaction) (string, error) {
	// Stub returns a dummy ID.
	return "stub-id", nil
}

func (s *StubTransactionRepo) Show(ctx context.Context, id string) (domain.Transaction, error) {
	// Stub returns zero value.
	return domain.Transaction{}, nil
}

func (s *StubTransactionRepo) Update(ctx context.Context, tx domain.Transaction) error {
	// Stub no‑op.
	return nil
}

func (s *StubTransactionRepo) Delete(ctx context.Context, id string) error {
	// Stub no‑op.
	return nil
}
