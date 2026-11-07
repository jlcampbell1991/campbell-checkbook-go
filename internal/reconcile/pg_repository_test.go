package reconcile

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
	mock.ExpectQuery("INSERT INTO import_jobs").WithArgs("PENDING", "path/to/file.csv", "", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	job := &ImportJob{FilePath: "path/to/file.csv"}
	if err := repo.Create(ctx, job); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if job.ID != 1 {
		t.Fatalf("expected ID 1, got %d", job.ID)
	}
	if job.Status != "PENDING" {
		t.Fatalf("expected status PENDING, got %s", job.Status)
	}

	// Expect select by id
	now := time.Now().UTC()
	mock.ExpectQuery("SELECT id, status, file_path, result_summary, created_at, updated_at FROM import_jobs WHERE id = \\$1").WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"id", "status", "file_path", "result_summary", "created_at", "updated_at"}).AddRow(1, "PENDING", "path/to/file.csv", "", now, now))

	got, err := repo.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if diff := cmp.Diff(job.FilePath, got.FilePath); diff != "" {
		t.Errorf("FilePath diff: %s", diff)
	}
	if got.ID != job.ID {
		t.Errorf("expected ID %d, got %d", job.ID, got.ID)
	}
	if got.Status != job.Status {
		t.Errorf("expected status %s, got %s", job.Status, got.Status)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
