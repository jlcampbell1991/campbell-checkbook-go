package handlers

import (
    "html/template"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

func TestAccountsNewForm(t *testing.T) {
    tmpl, err := template.New("test").Parse(`{{define "form.tmpl"}}<form></form>{{end}}`)
    if err != nil {
        t.Fatalf("template parse: %v", err)
    }
    repo := db.NewInMemoryAccountRepo()
    r := chi.NewRouter()
    RegisterAccountsRoutes(r, tmpl, repo)

    req := httptest.NewRequest("GET", "/accounts/new", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}

func TestAccountsCreate(t *testing.T) {
    tmpl, err := template.New("test").Parse(`{{define "form.tmpl"}}<form></form>{{end}}`)
    if err != nil {
        t.Fatalf("template parse: %v", err)
    }
    repo := db.NewInMemoryAccountRepo()
    r := chi.NewRouter()
    RegisterAccountsRoutes(r, tmpl, repo)

    form := "name=Checking"
    req := httptest.NewRequest("POST", "/accounts", strings.NewReader(form))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusSeeOther {
        t.Fatalf("expected redirect status, got %d", w.Code)
    }
    if len(repo.List()) != 1 {
        t.Fatalf("expected 1 account, got %d", len(repo.List()))
    }
    if repo.List()[0].Name != "Checking" {
        t.Fatalf("account name mismatch: %s", repo.List()[0].Name)
    }
}
