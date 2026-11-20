package db

import (
    "testing"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

func TestStubTransactionRepoMethods(t *testing.T) {
    repo := NewStubTransactionRepo()
    if _, err := repo.Create(nil, domain.Transaction{}); err != nil {
        t.Fatalf("Create returned error: %v", err)
    }
    if _, err := repo.Show(nil, "id"); err != nil {
        t.Fatalf("Show returned error: %v", err)
    }
    if err := repo.Update(nil, domain.Transaction{}); err != nil {
        t.Fatalf("Update returned error: %v", err)
    }
    if err := repo.Delete(nil, "id"); err != nil {
        t.Fatalf("Delete returned error: %v", err)
    }
}
