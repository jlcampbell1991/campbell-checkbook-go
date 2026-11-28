package db

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"testing"
)

func TestStubTransactionRepo(t *testing.T) {
	repo := NewStubTransactionRepo()
	// Test Create
	id, err := repo.Create(context.TODO(), domain.Transaction{})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if id == "" {
		t.Fatalf("Create returned empty ID")
	}
	// Test Show
	_, err = repo.Show(context.TODO(), id)
	if err != nil {
		t.Fatalf("Show returned error: %v", err)
	}
	// Test Update
	if err := repo.Update(context.TODO(), domain.Transaction{}); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	// Test Delete
	if err := repo.Delete(context.TODO(), id); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
}
