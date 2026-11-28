package db

import (
    "context"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

type PayDayDepositRepository interface {
    GetAllByPayDay(ctx context.Context, payDayId domain.PayDayId) ([]domain.PayDayDeposit, error)
    Upsert(ctx context.Context, dep domain.PayDayDeposit, payDayId domain.PayDayId) error
    GetById(ctx context.Context, id string) (domain.PayDayDeposit, error)
    Destroy(ctx context.Context, id string) error
}

type StubPayDayDepositRepo struct{}

func NewStubPayDayDepositRepo() PayDayDepositRepository { return &StubPayDayDepositRepo{} }

func (s *StubPayDayDepositRepo) GetAllByPayDay(ctx context.Context, payDayId domain.PayDayId) ([]domain.PayDayDeposit, error) {
    return []domain.PayDayDeposit{}, nil
}
func (s *StubPayDayDepositRepo) Upsert(ctx context.Context, dep domain.PayDayDeposit, payDayId domain.PayDayId) error { return nil }
func (s *StubPayDayDepositRepo) GetById(ctx context.Context, id string) (domain.PayDayDeposit, error) {
    return domain.PayDayDeposit{}, nil
}
func (s *StubPayDayDepositRepo) Destroy(ctx context.Context, id string) error { return nil }
