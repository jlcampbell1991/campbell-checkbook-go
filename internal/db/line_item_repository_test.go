package db

import (
    "context"
    "testing"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

func TestStubLineItemRepo_Insert(t *testing.T) {
    repo := NewStubLineItemRepo()
    li := domain.LineItem{ID: "", AccountID: "acct-1", Name: "Test", Balance: 123.45}
    id, err := repo.Insert(context.Background(), li)
    if err != nil {
        t.Fatalf("Insert returned error: %v", err)
    }
    if id == "" {
        t.Fatalf("Insert returned empty ID")
    }
    // Ensure placeholder ID is used when none provided
    if id != "stub-id" {
        t.Fatalf("expected stub-id, got %s", id)
    }
}
