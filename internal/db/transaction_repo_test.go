package db

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"testing"
)

func TestStubTransactionRepository_Create(t *testing.T) {
	repo := &StubTransactionRepository{}
	tx := domain.Transaction{Memo: "test"}
	id, err := repo.Create(context.Background(), tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatalf("expected non‑empty id")
	}
}

func TestStubTransactionRepository_Show(t *testing.T) {
	repo := &StubTransactionRepository{}
	tx, err := repo.Show(context.Background(), "some-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.ID != "some-id" {
		t.Fatalf("expected ID 'some-id', got %s", tx.ID)
	}
}
