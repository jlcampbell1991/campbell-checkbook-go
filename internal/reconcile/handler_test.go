package reconcile

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type memRepo struct {
	data map[int64]*ImportJob
	next int64
}

func newMemRepo() *memRepo {
	return &memRepo{data: make(map[int64]*ImportJob), next: 1}
}

func (m *memRepo) Create(ctx context.Context, job *ImportJob) error {
	job.ID = m.next
	m.next++
	now := time.Now().UTC()
	job.Status = "PENDING"
	job.CreatedAt = now
	job.UpdatedAt = now
	m.data[job.ID] = job
	return nil
}

func (m *memRepo) GetByID(ctx context.Context, id int64) (*ImportJob, error) {
	if j, ok := m.data[id]; ok {
		return j, nil
	}
	return nil, nil
}

func (m *memRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	if j, ok := m.data[id]; ok {
		j.Status = status
		j.UpdatedAt = time.Now().UTC()
		return nil
	}
	return nil
}

func (m *memRepo) List(ctx context.Context) ([]*ImportJob, error) {
	var list []*ImportJob
	for _, j := range m.data {
		list = append(list, j)
	}
	return list, nil
}

func TestReconcileHandlers(t *testing.T) {
	repo := newMemRepo()
	r := chi.NewRouter()
	RegisterRoutes(r, repo)
	server := httptest.NewServer(r)
	defer server.Close()

	// Create job
	createPayload := CreateImportJobRequest{FilePath: "test.csv"}
	body, _ := json.Marshal(createPayload)
	resp, err := http.Post(server.URL+"/reconcile/jobs", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var created ImportJob
	_ = json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	if created.ID == 0 {
		t.Fatalf("expected non-zero ID")
	}
	if created.Status != "PENDING" {
		t.Fatalf("expected status PENDING, got %s", created.Status)
	}

	// Get job
	getResp, err := http.Get(server.URL + "/reconcile/jobs/" + strconv.FormatInt(created.ID, 10))
	if err != nil {
		t.Fatalf("GET error: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", getResp.StatusCode)
	}
	var fetched ImportJob
	_ = json.NewDecoder(getResp.Body).Decode(&fetched)
	getResp.Body.Close()
	if fetched.ID != created.ID {
		t.Fatalf("fetched ID mismatch")
	}
	if fetched.FilePath != "test.csv" {
		t.Fatalf("expected file_path test.csv, got %s", fetched.FilePath)
	}

	// List jobs
	listResp, err := http.Get(server.URL + "/reconcile/jobs")
	if err != nil {
		t.Fatalf("LIST error: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("LIST expected 200, got %d", listResp.StatusCode)
	}
	var list []ImportJob
	_ = json.NewDecoder(listResp.Body).Decode(&list)
	listResp.Body.Close()
	if len(list) != 1 {
		t.Fatalf("expected list length 1, got %d", len(list))
	}
}
