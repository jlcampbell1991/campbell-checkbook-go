package accounts

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestPGRepository_CreateAndGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := NewPGRepository(db)
	ctx := context.Background()

	// Expect insert returning id
	mock.ExpectQuery("INSERT INTO accounts").WithArgs("Alice", int64(1000), "USD", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	acct := &Account{Name: "Alice", Balance: 1000, Currency: "USD"}
	if err := repo.Create(ctx, acct); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if acct.ID != 1 {
		t.Fatalf("expected ID 1, got %d", acct.ID)
	}

	// Expect select by id
	now := time.Now().UTC()
	mock.ExpectQuery("SELECT id, name, balance, currency, created_at, updated_at FROM accounts WHERE id = \\$1").WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "balance", "currency", "created_at", "updated_at"}).AddRow(1, "Alice", int64(1000), "USD", now, now))

	got, err := repo.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if diff := cmp.Diff(acct.Name, got.Name); diff != "" {
		t.Errorf("Name diff: %s", diff)
	}
	if got.ID != acct.ID {
		t.Errorf("expected ID %d, got %d", acct.ID, got.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
