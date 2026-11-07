package transactions

import (
	"bytes"
	"context"
	"encoding/csv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPGRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := NewPGRepository(db)
	ctx := context.Background()
	tx := &Transaction{AccountID: 1, LineItemID: 2, Amount: 1000, Currency: "USD", Date: time.Now()}

	// Expect query with returning id
	mock.ExpectQuery("INSERT INTO transactions").WithArgs(tx.AccountID, tx.LineItemID, tx.Amount, tx.Currency, tx.Date, tx.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	if err := repo.Create(ctx, tx); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if tx.ID != 10 {
		t.Fatalf("expected ID 10, got %d", tx.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPGRepository_ImportCSV(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := NewPGRepository(db)
	ctx := context.Background()

	// Prepare CSV data
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	header := []string{"account_id", "line_item_id", "amount", "currency", "date", "description"}
	_ = w.Write(header)
	_ = w.Write([]string{"1", "2", "1000", "USD", "2023-01-01", "Test desc"})
	_ = w.Write([]string{"3", "4", "2000", "EUR", "2023-02-02", "Another"})
	w.Flush()

	// Expectations: Begin, Prepare, Exec twice, Commit
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO transactions").ExpectExec().WithArgs(1, 2, int64(1000), "USD", sqlmock.AnyArg(), "Test desc", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactions").WithArgs(3, 4, int64(2000), "EUR", sqlmock.AnyArg(), "Another", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	count, err := repo.ImportCSV(ctx, bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ImportCSV error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 imported rows, got %d", count)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
