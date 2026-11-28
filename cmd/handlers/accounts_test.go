package handlers

import (
	"context"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockRepo implements db.AccountRepository for testing.
type mockRepo struct {
	inserted domain.Account
	called   bool
}

func (m *mockRepo) Insert(ctx context.Context, a domain.Account) (string, error) {
	m.inserted = a
	m.called = true
	return "mock-id", nil
}

func TestAccountNewHandler(t *testing.T) {
	tmpl, err := template.ParseFiles(
		"../../views/layouts/application.tmpl",
		"../../views/layouts/menu.tmpl",
		"../../views/accounts/add.tmpl",
	)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}
	h := &accountService{repo: nil, tmpl: tmpl}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/accounts/new", nil)
	h.new(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<form") || !strings.Contains(body, "action=\"/accounts\"") {
		t.Fatalf("response does not contain expected form: %s", body)
	}
}

func TestAccountCreateHandler(t *testing.T) {
	tmpl, err := template.ParseFiles(
		"../../views/layouts/application.tmpl",
		"../../views/layouts/menu.tmpl",
		"../../views/accounts/add.tmpl",
	)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}
	mr := &mockRepo{}
	h := &accountService{repo: mr, tmpl: tmpl}
	form := "name=TestAccount"
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.create(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status 303, got %d", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/" {
		t.Fatalf("expected Location '/' got %s", loc)
	}
	if !mr.called {
		t.Fatalf("repository Insert was not called")
	}
	if mr.inserted.Name != "TestAccount" {
		t.Fatalf("expected account name 'TestAccount', got %s", mr.inserted.Name)
	}
}
