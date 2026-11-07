package reconcile

import "context"

type Repository interface {
	Create(ctx context.Context, job *ImportJob) error
	GetByID(ctx context.Context, id int64) (*ImportJob, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	List(ctx context.Context) ([]*ImportJob, error)
}
