package transactions

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockRepo struct {
	imported int
	err      error
}

func (m *mockRepo) Create(ctx context.Context, tx *Transaction) error {
	return nil
}
func (m *mockRepo) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	return nil, nil
}
func (m *mockRepo) Update(ctx context.Context, tx *Transaction) error {
	return nil
}
func (m *mockRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
func (m *mockRepo) List(ctx context.Context) ([]*Transaction, error) {
	return nil, nil
}
func (m *mockRepo) ImportCSV(ctx context.Context, r io.Reader) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	// simple count lines ignoring header
	data, _ := io.ReadAll(r)
	lines := bytes.Count(data, []byte("\n"))
	// subtract header line
	if lines > 0 {
		lines--
	}
	m.imported = lines
	return lines, nil
}

func TestImportCSVHandler(t *testing.T) {
	repo := &mockRepo{}
	handler := importCSVHandler(repo)

	// Build multipart request with a small CSV
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	csvContent := "account_id,line_item_id,amount,currency,date,description\n1,2,1000,USD,2023-01-01,desc\n"
	part.Write([]byte(csvContent))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/transactions/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	// Expect JSON with imported count 1
	expected := `{"imported":1}`
	if strings.TrimSpace(w.Body.String()) != expected {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
	if repo.imported != 1 {
		t.Fatalf("repo imported count mismatch: %d", repo.imported)
	}
}
