package handlers

import (
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"
    "text/template"

    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

func TestTransactionNewHandler(t *testing.T) {
    tmpl := template.Must(template.New("test").Parse(`{{define \"transactions_add.tmpl\"}}<form></form>{{end}}`))
    if tmpl == nil {
        t.Fatalf("failed to parse template")
    }
    // Use stub repo (no need for real implementation)
    repo := db.TransactionRepo
    router := TransactionRoutes(tmpl, repo)
    req := httptest.NewRequest(http.MethodGet, "/new?account_id=acc1&line_item_id=li1", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}

func TestTransactionCreateHandler(t *testing.T) {
    tmpl := template.New("test").Parse(`{{define \"transactions_add.tmpl\"}}<form></form>{{end}}`)
    repo := db.TransactionRepo
    router := TransactionRoutes(tmpl, repo)
    form := url.Values{}
    form.Set("date", "2023-01-01")
    form.Set("memo", "Test memo")
    form.Set("amount", "12.34")
    form.Set("category", "Misc")
    form.Set("account_id", "acc1")
    form.Set("line_item_id", "li1")
    req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    if w.Code != http.StatusSeeOther {
        t.Fatalf("expected redirect status 303, got %d", w.Code)
    }
}
