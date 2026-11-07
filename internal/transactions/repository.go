package transactions

import (
	"context"
	"io"
)

// Repository defines data‑access methods for Transaction entities.
type Repository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	Update(ctx context.Context, tx *Transaction) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Transaction, error)
	// ImportCSV reads CSV data from r and creates transactions. Returns number of records imported.
	ImportCSV(ctx context.Context, r io.Reader) (int, error)
}
