package reconcile

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PGRepository implements Repository using a *sql.DB with pgx driver.
type PGRepository struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) *PGRepository {
	return &PGRepository{db: db}
}

func (r *PGRepository) Create(ctx context.Context, job *ImportJob) error {
	if job == nil {
		return fmt.Errorf("job is nil")
	}
	query := `INSERT INTO import_jobs (status, file_path, result_summary, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	now := time.Now().UTC()
	job.Status = "PENDING"
	job.CreatedAt = now
	job.UpdatedAt = now
	err := r.db.QueryRowContext(ctx, query, job.Status, job.FilePath, job.ResultSummary, now, now).Scan(&job.ID)
	if err != nil {
		return fmt.Errorf("insert import job: %w", err)
	}
	return nil
}

func (r *PGRepository) GetByID(ctx context.Context, id int64) (*ImportJob, error) {
	query := `SELECT id, status, file_path, result_summary, created_at, updated_at FROM import_jobs WHERE id = $1`
	job := &ImportJob{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&job.ID, &job.Status, &job.FilePath, &job.ResultSummary, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select import job: %w", err)
	}
	return job, nil
}

func (r *PGRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE import_jobs SET status = $1, updated_at = $2 WHERE id = $3`
	now := time.Now().UTC()
	res, err := r.db.ExecContext(ctx, query, status, now, id)
	if err != nil {
		return fmt.Errorf("update import job status: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

func (r *PGRepository) List(ctx context.Context) ([]*ImportJob, error) {
	query := `SELECT id, status, file_path, result_summary, created_at, updated_at FROM import_jobs ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list import jobs: %w", err)
	}
	defer rows.Close()
	var jobs []*ImportJob
	for rows.Next() {
		job := &ImportJob{}
		if err := rows.Scan(&job.ID, &job.Status, &job.FilePath, &job.ResultSummary, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan import job: %w", err)
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return jobs, nil
}
