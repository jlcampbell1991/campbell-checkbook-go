package lineitems

import "context"

type Repository interface {
	Create(ctx context.Context, li *LineItem) error
	GetByID(ctx context.Context, id int64) (*LineItem, error)
	Update(ctx context.Context, li *LineItem) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*LineItem, error)
}
