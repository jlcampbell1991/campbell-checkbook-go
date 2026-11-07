package lineitems

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PGRepository implements Repository using a *sql.DB with pgx driver.
type PGRepository struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) *PGRepository {
	return &PGRepository{db: db}
}

func (r *PGRepository) Create(ctx context.Context, li *LineItem) error {
	if li == nil {
		return fmt.Errorf("line item is nil")
	}
	query := `INSERT INTO line_items (account_id, description, amount, currency, date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	now := time.Now().UTC()
	li.CreatedAt = now
	li.UpdatedAt = now
	err := r.db.QueryRowContext(ctx, query, li.AccountID, li.Description, li.Amount, li.Currency, li.Date, now, now).Scan(&li.ID)
	if err != nil {
		return fmt.Errorf("insert line item: %w", err)
	}
	return nil
}

func (r *PGRepository) GetByID(ctx context.Context, id int64) (*LineItem, error) {
	query := `SELECT id, account_id, description, amount, currency, date, created_at, updated_at FROM line_items WHERE id = $1`
	li := &LineItem{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&li.ID, &li.AccountID, &li.Description, &li.Amount, &li.Currency, &li.Date, &li.CreatedAt, &li.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select line item: %w", err)
	}
	return li, nil
}

func (r *PGRepository) Update(ctx context.Context, li *LineItem) error {
	if li == nil {
		return fmt.Errorf("line item is nil")
	}
	query := `UPDATE line_items SET account_id = $1, description = $2, amount = $3, currency = $4, date = $5, updated_at = $6 WHERE id = $7`
	now := time.Now().UTC()
	li.UpdatedAt = now
	res, err := r.db.ExecContext(ctx, query, li.AccountID, li.Description, li.Amount, li.Currency, li.Date, now, li.ID)
	if err != nil {
		return fmt.Errorf("update line item: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

func (r *PGRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM line_items WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete line item: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no rows deleted")
	}
	return nil
}

func (r *PGRepository) List(ctx context.Context) ([]*LineItem, error) {
	query := `SELECT id, account_id, description, amount, currency, date, created_at, updated_at FROM line_items ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list line items: %w", err)
	}
	defer rows.Close()
	var items []*LineItem
	for rows.Next() {
		li := &LineItem{}
		if err := rows.Scan(&li.ID, &li.AccountID, &li.Description, &li.Amount, &li.Currency, &li.Date, &li.CreatedAt, &li.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan line item: %w", err)
		}
		items = append(items, li)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return items, nil
}
