package db

import (
    "context"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

type AccountRepository interface {
    // Minimal methods used by PayDayRoutes (if any). For now stub only.
    GetAll(ctx context.Context) ([]domain.Account, error)
    GetById(ctx context.Context, id string) (domain.Account, error)
}

type StubAccountRepo struct{}

func NewStubAccountRepo() AccountRepository { return &StubAccountRepo{} }

func (s *StubAccountRepo) GetAll(ctx context.Context) ([]domain.Account, error) { return []domain.Account{}, nil }
func (s *StubAccountRepo) GetById(ctx context.Context, id string) (domain.Account, error) { return domain.Account{}, nil }
