package accounts

import "context"

type Repository interface {
	Create(ctx context.Context, acct *Account) error
	GetByID(ctx context.Context, id int64) (*Account, error)
	Update(ctx context.Context, acct *Account) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Account, error)
}
