package utils

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
)

func TestScanNullableString(t *testing.T) {
	var dest *string
	// non‑null []byte
	if err := ScanNullableString([]byte("hello"), &dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest == nil || *dest != "hello" {
		t.Fatalf("expected 'hello', got %v", dest)
	}
	// null value
	if err := ScanNullableString(nil, &dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest != nil {
		t.Fatalf("expected nil, got %v", dest)
	}
	// unsupported type
	if err := ScanNullableString(123, &dest); err == nil {
		t.Fatalf("expected error for unsupported type")
	}
}

func TestExecWithTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	// Expect transaction begin, exec, commit.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO test (name) VALUES (?)")).WithArgs("bob").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = ExecWithTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO test (name) VALUES (?)", "bob")
		return err
	})
	if err != nil {
		t.Fatalf("ExecWithTx returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
