package paydays

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPGRepository_CreatePayDay(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := NewPGRepository(db)
	ctx := context.Background()
	pd := &PayDay{Date: time.Now(), Amount: 1000, Currency: "USD"}

	mock.ExpectQuery("INSERT INTO pay_days").
		WithArgs(pd.Date, pd.Amount, pd.Currency, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	if err := repo.CreatePayDay(ctx, pd); err != nil {
		t.Fatalf("CreatePayDay error: %v", err)
	}
	if pd.ID != 1 {
		t.Fatalf("expected ID 1, got %d", pd.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
