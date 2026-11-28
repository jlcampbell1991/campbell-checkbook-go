package db

import (
    "context"
    "testing"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

func TestStubTransactionRepoCreate(t *testing.T) {
    repo := &StubTransactionRepo{}
    tx := domain.Transaction{Memo: "test"}
    id, err := repo.Create(context.Background(), tx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if id == "" {
        t.Fatalf("expected non‑empty id")
    }
}
