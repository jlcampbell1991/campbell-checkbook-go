package db

import (
	"context"
	"database/sql"
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
func (r *AccountRepository) Index(ctx context.Context) ([]domain.Account, error) {
    // Query accounts with aggregated balance and budget.
    rows, err := r.Pool.Query(ctx, "SELECT accounts.name, SUM(items.balance) AS balance, SUM(deps.amount) AS budget, accounts.id FROM checkbook_accounts accounts LEFT JOIN checkbook_line_items items ON (items.account_id = accounts.id AND items.active = true) LEFT JOIN checkbook_pay_day_deposits deps ON (deps.line_item_id = items.id) GROUP BY accounts.name, accounts.id ORDER BY accounts.name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var accounts []domain.Account
    for rows.Next() {
        var name string
        var balance sql.NullFloat64
        var budget sql.NullFloat64
        var id string
        if err := rows.Scan(&name, &balance, &budget, &id); err != nil {
            return nil, err
        }

        // Fetch line items for this account.
        liRows, err := r.Pool.Query(ctx, "SELECT items.name, SUM(deps.amount) AS budget, items.balance, items.active, items.account_id, items.id FROM checkbook_line_items items LEFT JOIN checkbook_pay_day_deposits deps ON deps.line_item_id = items.id WHERE items.active = true AND items.account_id = $1 GROUP BY items.name, items.balance, items.active, items.account_id, items.id ORDER BY items.name", id)
        if err != nil {
            return nil, err
        }
        var lineItems []domain.LineItem
        for liRows.Next() {
            var li domain.LineItem
            var liBal sql.NullFloat64
            var liBud sql.NullFloat64
            var active bool
            if err := liRows.Scan(&li.Name, &liBud, &liBal, &active, &li.AccountID, &li.ID); err != nil {
                liRows.Close()
                return nil, err
            }
            if liBal.Valid {
                li.Balance = &liBal.Float64
            }
            if liBud.Valid {
                li.Budget = &liBud.Float64
            }
            li.Active = active
            lineItems = append(lineItems, li)
        }
        liRows.Close()

        acct := domain.Account{
            ID:        id,
            Name:      name,
            LineItems: lineItems,
        }
        if balance.Valid {
            acct.Balance = &balance.Float64
        }
        if budget.Valid {
            acct.Budget = &budget.Float64
        }
        accounts = append(accounts, acct)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return accounts, nil
}
