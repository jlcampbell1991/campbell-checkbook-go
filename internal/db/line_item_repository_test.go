package db

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"testing"
)

func TestStubLineItemRepository_Insert(t *testing.T) {
	repo := &StubLineItemRepository{}
	li := domain.LineItem{ID: "", AccountID: "acc1", Name: "Test", Balance: 10.5}
	id, err := repo.Insert(context.Background(), li)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatalf("expected non‑empty id, got empty")
	}
}
