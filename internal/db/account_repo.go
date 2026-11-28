package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

type AccountRepository struct {
	Pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{Pool: pool}
}

// Index returns a list of accounts with their line items.
// TODO: implement actual SQL query joining accounts and line items.
func (r *AccountRepository) Index(ctx context.Context) ([]domain.Account, error) {
	// Placeholder implementation returning empty slice.
	return []domain.Account{}, nil
}
