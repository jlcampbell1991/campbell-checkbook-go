package paydays

import "context"

// Repository defines data‑access methods for PayDay and PayDayDeposit entities.
type Repository interface {
	// PayDay methods
	CreatePayDay(ctx context.Context, pd *PayDay) error
	GetPayDayByID(ctx context.Context, id int64) (*PayDay, error)
	UpdatePayDay(ctx context.Context, pd *PayDay) error
	DeletePayDay(ctx context.Context, id int64) error
	ListPayDays(ctx context.Context) ([]*PayDay, error)

	// PayDayDeposit methods
	CreateDeposit(ctx context.Context, d *PayDayDeposit) error
	GetDepositByID(ctx context.Context, id int64) (*PayDayDeposit, error)
	UpdateDeposit(ctx context.Context, d *PayDayDeposit) error
	DeleteDeposit(ctx context.Context, id int64) error
	ListDepositsByPayDay(ctx context.Context, paydayID int64) ([]*PayDayDeposit, error)
}
