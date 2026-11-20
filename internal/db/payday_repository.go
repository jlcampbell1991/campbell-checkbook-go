package db

import (
    "context"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// PayDayRepository defines data access methods needed by the pay‑day handlers.
type PayDayRepository interface {
    GetAll(ctx context.Context) ([]domain.PayDay, error)
    GetAllInactive(ctx context.Context) ([]domain.PayDay, error)
    GetById(ctx context.Context, id string) (domain.PayDay, error)
    Insert(ctx context.Context, req *domain.PayDayReq) (*domain.PayDay, error)
    Update(ctx context.Context, req *domain.PayDayReq, id string) (*domain.PayDay, error)
    Deactivate(ctx context.Context, id string) error
    Reactivate(ctx context.Context, id string) error
}

// StubPayDayRepo provides a minimal in‑memory implementation for compilation.
type StubPayDayRepo struct{}

func NewStubPayDayRepo() PayDayRepository { return &StubPayDayRepo{} }

func (s *StubPayDayRepo) GetAll(ctx context.Context) ([]domain.PayDay, error) { return nil, nil }
func (s *StubPayDayRepo) GetAllInactive(ctx context.Context) ([]domain.PayDay, error) { return nil, nil }
func (s *StubPayDayRepo) GetById(ctx context.Context, id string) (domain.PayDay, error) {
    return domain.PayDay{Id: domain.PayDayId{Value: id}, Name: "Stub PayDay", Total: 0, AccountId: ""}, nil
}
func (s *StubPayDayRepo) Insert(ctx context.Context, req *domain.PayDayReq) (*domain.PayDay, error) {
    pd := &domain.PayDay{Id: domain.PayDayId{Value: "stub-id"}, Name: req.Name, Total: req.Total, AccountId: req.AccountId}
    return pd, nil
}
func (s *StubPayDayRepo) Update(ctx context.Context, req *domain.PayDayReq, id string) (*domain.PayDay, error) {
    pd := &domain.PayDay{Id: domain.PayDayId{Value: id}, Name: req.Name, Total: req.Total, AccountId: req.AccountId}
    return pd, nil
}
func (s *StubPayDayRepo) Deactivate(ctx context.Context, id string) error { return nil }
func (s *StubPayDayRepo) Reactivate(ctx context.Context, id string) error { return nil }
