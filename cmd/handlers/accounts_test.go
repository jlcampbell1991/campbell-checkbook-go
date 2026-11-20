package handlers

import (
	"bytes"
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockRepo implements db.AccountRepository for testing.
type mockRepo struct {
	inserted     domain.Account
	insertCalled bool
}

func (m *mockRepo) Index(ctx context.Context) ([]domain.Account, error) { return nil, nil }
func (m *mockRepo) Show(ctx context.Context, id string) (domain.Account, error) {
	return domain.Account{}, nil
}
func (m *mockRepo) Insert(ctx context.Context, a domain.Account) (string, error) {
	m.inserted = a
	m.insertCalled = true
	return "1", nil
}
func (m *mockRepo) Update(ctx context.Context, a domain.Account, id string) error { return nil }
func (m *mockRepo) Destroy(ctx context.Context, id string) error                  { return nil }
func (m *mockRepo) GetTransactions(ctx context.Context, id string, page int) ([]domain.Transaction, error) {
	return nil, nil
}
func (m *mockRepo) GetDislocated(ctx context.Context, id string) ([]domain.Transaction, error) {
	return nil, nil
}
func (m *mockRepo) RelocateAll(ctx context.Context, id string) error { return nil }

func TestAccountNewHandler(t *testing.T) {
	// Prepare a minimal template set containing accounts_add.tmpl.
	tmpl := template.New("test")
	tmpl, err := tmpl.Parse(`{{define "accounts_add.tmpl"}}<form action="/accounts" method="post">Name: <input name="name"></form>{{end}}`)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}
	h := &accountService{repo: &mockRepo{}, tmpl: tmpl}

	req := httptest.NewRequest(http.MethodGet, "/accounts/new", nil)
	w := httptest.NewRecorder()
	h.new(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte("<form")) || !bytes.Contains([]byte(body), []byte("action=\"/accounts\"")) {
		t.Fatalf("response does not contain expected form: %s", body)
	}
}

func TestAccountCreateHandler(t *testing.T) {
	tmpl := template.New("test") // no template needed for redirect path
	repo := &mockRepo{}
	h := &accountService{repo: repo, tmpl: tmpl}

	form := "name=TestAccount"
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.create(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if loc != "/" {
		t.Fatalf("expected redirect to '/', got %s", loc)
	}
	if !repo.insertCalled {
		t.Fatalf("repo Insert was not called")
	}
	if repo.inserted.Name != "TestAccount" {
		t.Fatalf("expected inserted name 'TestAccount', got '%s'", repo.inserted.Name)
	}
}
