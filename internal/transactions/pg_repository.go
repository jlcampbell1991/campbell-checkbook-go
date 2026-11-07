package transactions

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

// PGRepository implements Repository using a *sql.DB with pgx driver.
type PGRepository struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) *PGRepository {
	return &PGRepository{db: db}
}

func (r *PGRepository) Create(ctx context.Context, tx *Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	query := `INSERT INTO transactions (account_id, line_item_id, amount, currency, date, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	now := time.Now().UTC()
	tx.CreatedAt = now
	tx.UpdatedAt = now
	err := r.db.QueryRowContext(ctx, query, tx.AccountID, tx.LineItemID, tx.Amount, tx.Currency, tx.Date, tx.Description, now, now).Scan(&tx.ID)
	if err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}
	return nil
}

func (r *PGRepository) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	query := `SELECT id, account_id, line_item_id, amount, currency, date, description, created_at, updated_at FROM transactions WHERE id = $1`
	tx := &Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&tx.ID, &tx.AccountID, &tx.LineItemID, &tx.Amount, &tx.Currency, &tx.Date, &tx.Description, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select transaction: %w", err)
	}
	return tx, nil
}

func (r *PGRepository) Update(ctx context.Context, tx *Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	query := `UPDATE transactions SET account_id = $1, line_item_id = $2, amount = $3, currency = $4, date = $5, description = $6, updated_at = $7 WHERE id = $8`
	now := time.Now().UTC()
	tx.UpdatedAt = now
	res, err := r.db.ExecContext(ctx, query, tx.AccountID, tx.LineItemID, tx.Amount, tx.Currency, tx.Date, tx.Description, now, tx.ID)
	if err != nil {
		return fmt.Errorf("update transaction: %w", err)
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
	query := `DELETE FROM transactions WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete transaction: %w", err)
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

func (r *PGRepository) List(ctx context.Context) ([]*Transaction, error) {
	query := `SELECT id, account_id, line_item_id, amount, currency, date, description, created_at, updated_at FROM transactions ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()
	var list []*Transaction
	for rows.Next() {
		tx := &Transaction{}
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.LineItemID, &tx.Amount, &tx.Currency, &tx.Date, &tx.Description, &tx.CreatedAt, &tx.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		list = append(list, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return list, nil
}

// ImportCSV parses CSV data and creates transactions in a single transaction.
func (r *PGRepository) ImportCSV(ctx context.Context, reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	// Expect header row.
	records, err := csvReader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("read csv: %w", err)
	}
	if len(records) < 1 {
		return 0, fmt.Errorf("csv empty")
	}
	header := records[0]
	// Expected columns (case‑insensitive): account_id,line_item_id,amount,currency,date,description
	expected := []string{"account_id", "line_item_id", "amount", "currency", "date", "description"}
	if len(header) != len(expected) {
		return 0, fmt.Errorf("invalid csv header, expected %d columns", len(expected))
	}
	for i, col := range expected {
		if header[i] != col {
			return 0, fmt.Errorf("unexpected column %s, expected %s", header[i], col)
		}
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO transactions (account_id, line_item_id, amount, currency, date, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		return 0, fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()
	count := 0
	now := time.Now().UTC()
	for _, rec := range records[1:] {
		if len(rec) != len(expected) {
			return count, fmt.Errorf("invalid record length")
		}
		accountID, err := strconv.ParseInt(rec[0], 10, 64)
		if err != nil {
			return count, fmt.Errorf("parse account_id: %w", err)
		}
		lineItemID, err := strconv.ParseInt(rec[1], 10, 64)
		if err != nil {
			return count, fmt.Errorf("parse line_item_id: %w", err)
		}
		amount, err := strconv.ParseInt(rec[2], 10, 64)
		if err != nil {
			return count, fmt.Errorf("parse amount: %w", err)
		}
		currency := rec[3]
		date, err := time.Parse("2006-01-02", rec[4])
		if err != nil {
			return count, fmt.Errorf("parse date: %w", err)
		}
		description := rec[5]
		// Simple validation
		if accountID <= 0 || lineItemID <= 0 || amount < 0 || currency == "" || description == "" {
			return count, fmt.Errorf("validation error on record %d", count+1)
		}
		_, err = stmt.ExecContext(ctx, accountID, lineItemID, amount, currency, date, description, now, now)
		if err != nil {
			return count, fmt.Errorf("exec insert: %w", err)
		}
		count++
	}
	if err := tx.Commit(); err != nil {
		return count, fmt.Errorf("commit tx: %w", err)
	}
	return count, nil
}
