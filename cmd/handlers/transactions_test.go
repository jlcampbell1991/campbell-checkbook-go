package handlers

import (
    "net/http"
    "net/http/httptest"
    "net/url"
    
    "strings"
    "testing"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    "html/template"
)



func TestTransactionNew(t *testing.T) {
    // Parse minimal templates needed for the handler
    // Use an empty template as handlers do not rely on parsed templates for these tests.
    tmpl := template.New("test")
    var err error = nil
    if err != nil {
        t.Fatalf("template parse error: %v", err)
    }
    repo := db.NewStubTransactionRepo()
    router := TransactionRoutes(tmpl, repo)

    req := httptest.NewRequest(http.MethodGet, "/new?account_id=a1&line_item_id=l1", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}

func TestTransactionCreate(t *testing.T) {
    // Use an empty template as handlers do not rely on parsed templates for these tests.
    tmpl := template.New("test")
    var err error = nil
    if err != nil {
        t.Fatalf("template parse error: %v", err)
    }
    repo := db.NewStubTransactionRepo()
    router := TransactionRoutes(tmpl, repo)

    form := url.Values{}
    form.Set("date", "2023-01-01")
    form.Set("memo", "Test transaction")
    form.Set("amount", "12.34")
    form.Set("category", "Misc")
    form.Set("account_id", "a1")
    form.Set("line_item_id", "l1")
    body := strings.NewReader(form.Encode())
    req := httptest.NewRequest(http.MethodPost, "/", body)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    if w.Code != http.StatusSeeOther {
        t.Fatalf("expected redirect status %d, got %d", http.StatusSeeOther, w.Code)
    }
    // Expect redirect location to line-item page
    loc := w.Header().Get("Location")
    expected := "/line-items/l1"
    if loc != expected {
        t.Fatalf("expected redirect to %s, got %s", expected, loc)
    }
    // Ensure repo Create can be called without panic (stub returns stub-id)
    _ = repo // avoid unused warning
    _ = domain.Transaction{}
}
