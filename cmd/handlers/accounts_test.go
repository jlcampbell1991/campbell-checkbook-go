package handlers

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "github.com/go-chi/chi/v5"
    "campbell-checkbook/internal/db"
)

func TestAccountFormRender(t *testing.T) {
    repo := db.NewMemAccountRepo()
    r := chi.NewRouter()
    RegisterAccounts(r, repo)

    req := httptest.NewRequest(http.MethodGet, "/accounts/new", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "<form") {
        t.Fatalf("response does not contain form")
    }
}

func TestAccountCreate(t *testing.T) {
    repo := db.NewMemAccountRepo()
    r := chi.NewRouter()
    RegisterAccounts(r, repo)

    form := strings.NewReader("name=TestAccount")
    req := httptest.NewRequest(http.MethodPost, "/accounts", form)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusSeeOther {
        t.Fatalf("expected redirect (303), got %d", w.Code)
    }
    // Verify account stored
    accounts, _ := repo.Index()
    if len(accounts) != 1 || accounts[0].Name != "TestAccount" {
        t.Fatalf("account not stored correctly: %+v", accounts)
    }
}
