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

// mockAccountRepo implements AccountRepository for testing.
type mockAccountRepo struct {
	inserted domain.Account
	err      error
}

func (m *mockAccountRepo) Insert(ctx context.Context, a domain.Account) (string, error) {
	m.inserted = a
	if m.err != nil {
		return "", m.err
	}
	return "mock-id", nil
}

func TestAccountNewHandler(t *testing.T) {
	// Minimal templates: application and accounts_add.
	tmpl := template.New("base")
	tmpl, _ = tmpl.Parse(`{{define "application.tmpl"}}{{end}}{{define "accounts_add.tmpl"}}{{template "application.tmpl" .}}<form></form>{{end}}`)
	h := &accountService{repo: &mockAccountRepo{}, tmpl: tmpl}

	req := httptest.NewRequest(http.MethodGet, "/accounts/new", nil)
	w := httptest.NewRecorder()
	h.new(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("<form")) {
		t.Fatalf("response does not contain form element: %s", w.Body.String())
	}
}

func TestAccountCreateHandlerSuccess(t *testing.T) {
	tmpl := template.New("base")
	tmpl, _ = tmpl.Parse(`{{define "application.tmpl"}}{{end}}{{define "accounts_add.tmpl"}}{{template "application.tmpl" .}}<form></form>{{end}}`)
	repo := &mockAccountRepo{}
	h := &accountService{repo: repo, tmpl: tmpl}

	body := bytes.NewBufferString("name=TestAccount")
	req := httptest.NewRequest(http.MethodPost, "/accounts", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.create(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/" {
		t.Fatalf("expected redirect to '/', got %s", loc)
	}
	if repo.inserted.Name != "TestAccount" {
		t.Fatalf("expected inserted name 'TestAccount', got %s", repo.inserted.Name)
	}
}

func TestAccountCreateHandlerMissingName(t *testing.T) {
	tmpl := template.New("base")
	tmpl, _ = tmpl.Parse(`{{define "application.tmpl"}}{{end}}{{define "accounts_add.tmpl"}}{{template "application.tmpl" .}}<form></form>{{end}}`)
	h := &accountService{repo: &mockAccountRepo{}, tmpl: tmpl}

	body := bytes.NewBufferString("")
	req := httptest.NewRequest(http.MethodPost, "/accounts", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for missing name, got %d", http.StatusBadRequest, w.Code)
	}
}
