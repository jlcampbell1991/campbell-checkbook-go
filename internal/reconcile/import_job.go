package reconcile

import "time"

// ImportJob represents a background job for importing reconciliation data.
// Status values are simple strings; in a full implementation they could be an enum.
type ImportJob struct {
	ID            int64     `json:"id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	FilePath      string    `json:"file_path"`
	ResultSummary string    `json:"result_summary,omitempty"`
}
