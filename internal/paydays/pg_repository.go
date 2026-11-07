package paydays

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

// ---------- PayDay ----------
func (r *PGRepository) CreatePayDay(ctx context.Context, pd *PayDay) error {
	if pd == nil {
		return fmt.Errorf("payday is nil")
	}
	query := `INSERT INTO pay_days (date, amount, currency, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5) RETURNING id`
	now := time.Now().UTC()
	pd.CreatedAt = now
	pd.UpdatedAt = now
	err := r.db.QueryRowContext(ctx, query, pd.Date, pd.Amount, pd.Currency, now, now).Scan(&pd.ID)
	if err != nil {
		return fmt.Errorf("insert payday: %w", err)
	}
	return nil
}

func (r *PGRepository) GetPayDayByID(ctx context.Context, id int64) (*PayDay, error) {
	query := `SELECT id, date, amount, currency, created_at, updated_at FROM pay_days WHERE id = $1`
	pd := &PayDay{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&pd.ID, &pd.Date, &pd.Amount, &pd.Currency, &pd.CreatedAt, &pd.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select payday: %w", err)
	}
	return pd, nil
}

func (r *PGRepository) UpdatePayDay(ctx context.Context, pd *PayDay) error {
	if pd == nil {
		return fmt.Errorf("payday is nil")
	}
	query := `UPDATE pay_days SET date = $1, amount = $2, currency = $3, updated_at = $4 WHERE id = $5`
	now := time.Now().UTC()
	pd.UpdatedAt = now
	res, err := r.db.ExecContext(ctx, query, pd.Date, pd.Amount, pd.Currency, now, pd.ID)
	if err != nil {
		return fmt.Errorf("update payday: %w", err)
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

func (r *PGRepository) DeletePayDay(ctx context.Context, id int64) error {
	query := `DELETE FROM pay_days WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete payday: %w", err)
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

func (r *PGRepository) ListPayDays(ctx context.Context) ([]*PayDay, error) {
	query := `SELECT id, date, amount, currency, created_at, updated_at FROM pay_days ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list paydays: %w", err)
	}
	defer rows.Close()
	var list []*PayDay
	for rows.Next() {
		pd := &PayDay{}
		if err := rows.Scan(&pd.ID, &pd.Date, &pd.Amount, &pd.Currency, &pd.CreatedAt, &pd.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan payday: %w", err)
		}
		list = append(list, pd)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return list, nil
}

// ---------- PayDayDeposit ----------
func (r *PGRepository) CreateDeposit(ctx context.Context, d *PayDayDeposit) error {
	if d == nil {
		return fmt.Errorf("deposit is nil")
	}
	query := `INSERT INTO pay_day_deposits (payday_id, account_id, amount, currency, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	err := r.db.QueryRowContext(ctx, query, d.PayDayID, d.AccountID, d.Amount, d.Currency, now, now).Scan(&d.ID)
	if err != nil {
		return fmt.Errorf("insert deposit: %w", err)
	}
	return nil
}

func (r *PGRepository) GetDepositByID(ctx context.Context, id int64) (*PayDayDeposit, error) {
	query := `SELECT id, payday_id, account_id, amount, currency, created_at, updated_at FROM pay_day_deposits WHERE id = $1`
	d := &PayDayDeposit{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&d.ID, &d.PayDayID, &d.AccountID, &d.Amount, &d.Currency, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select deposit: %w", err)
	}
	return d, nil
}

func (r *PGRepository) UpdateDeposit(ctx context.Context, d *PayDayDeposit) error {
	if d == nil {
		return fmt.Errorf("deposit is nil")
	}
	query := `UPDATE pay_day_deposits SET payday_id = $1, account_id = $2, amount = $3, currency = $4, updated_at = $5 WHERE id = $6`
	now := time.Now().UTC()
	d.UpdatedAt = now
	res, err := r.db.ExecContext(ctx, query, d.PayDayID, d.AccountID, d.Amount, d.Currency, now, d.ID)
	if err != nil {
		return fmt.Errorf("update deposit: %w", err)
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

func (r *PGRepository) DeleteDeposit(ctx context.Context, id int64) error {
	query := `DELETE FROM pay_day_deposits WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete deposit: %w", err)
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

func (r *PGRepository) ListDepositsByPayDay(ctx context.Context, paydayID int64) ([]*PayDayDeposit, error) {
	query := `SELECT id, payday_id, account_id, amount, currency, created_at, updated_at FROM pay_day_deposits WHERE payday_id = $1 ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query, paydayID)
	if err != nil {
		return nil, fmt.Errorf("list deposits: %w", err)
	}
	defer rows.Close()
	var list []*PayDayDeposit
	for rows.Next() {
		d := &PayDayDeposit{}
		if err := rows.Scan(&d.ID, &d.PayDayID, &d.AccountID, &d.Amount, &d.Currency, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan deposit: %w", err)
		}
		list = append(list, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return list, nil
}
