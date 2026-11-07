package accounts

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

func (r *PGRepository) Create(ctx context.Context, acct *Account) error {
	if acct == nil {
		return fmt.Errorf("account is nil")
	}
	query := `INSERT INTO accounts (name, balance, currency, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	now := time.Now().UTC()
	acct.CreatedAt = now
	acct.UpdatedAt = now
	err := r.db.QueryRowContext(ctx, query, acct.Name, acct.Balance, acct.Currency, now, now).Scan(&acct.ID)
	if err != nil {
		return fmt.Errorf("insert account: %w", err)
	}
	return nil
}

func (r *PGRepository) GetByID(ctx context.Context, id int64) (*Account, error) {
	query := `SELECT id, name, balance, currency, created_at, updated_at FROM accounts WHERE id = $1`
	acct := &Account{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&acct.ID, &acct.Name, &acct.Balance, &acct.Currency, &acct.CreatedAt, &acct.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select account: %w", err)
	}
	return acct, nil
}

func (r *PGRepository) Update(ctx context.Context, acct *Account) error {
	if acct == nil {
		return fmt.Errorf("account is nil")
	}
	query := `UPDATE accounts SET name = $1, balance = $2, currency = $3, updated_at = $4 WHERE id = $5`
	now := time.Now().UTC()
	acct.UpdatedAt = now
	res, err := r.db.ExecContext(ctx, query, acct.Name, acct.Balance, acct.Currency, now, acct.ID)
	if err != nil {
		return fmt.Errorf("update account: %w", err)
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
	query := `DELETE FROM accounts WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete account: %w", err)
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

func (r *PGRepository) List(ctx context.Context) ([]*Account, error) {
	query := `SELECT id, name, balance, currency, created_at, updated_at FROM accounts ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()
	var accounts []*Account
	for rows.Next() {
		acct := &Account{}
		if err := rows.Scan(&acct.ID, &acct.Name, &acct.Balance, &acct.Currency, &acct.CreatedAt, &acct.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, acct)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return accounts, nil
}
