package handlers

import (
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestTransactionRoutes_New(t *testing.T) {
	tmpl, err := template.New("test").Parse(`{{define "transactions_add.tmpl"}}<form></form>{{end}}`)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}
	repo := db.NewStubTransactionRepo()
	r := TransactionRoutes(tmpl, repo)

	req := httptest.NewRequest(http.MethodGet, "/new?account_id=acc1&line_item_id=li1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestTransactionRoutes_Create_Success(t *testing.T) {
	tmpl := template.New("test")
	// No templates needed for create handler success path
	if tmpl == nil {
		t.Fatalf("failed to create test template")
	}
	repo := db.NewStubTransactionRepo()
	r := TransactionRoutes(tmpl, repo)

	form := url.Values{}
	form.Set("date", "2023-01-01")
	form.Set("memo", "Test transaction")
	form.Set("amount", "12.34")
	form.Set("category", "Misc")
	form.Set("account_id", "acc1")
	form.Set("line_item_id", "li1")
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 SeeOther, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	expected := "/line-items/li1"
	if loc != expected {
		t.Fatalf("expected redirect to %s, got %s", expected, loc)
	}
}

func TestTransactionRoutes_Create_MissingFields(t *testing.T) {
	tmpl, _ := template.ParseFiles("views/layouts/application.tmpl", "views/layouts/menu.tmpl")
	repo := db.NewStubTransactionRepo()
	r := TransactionRoutes(tmpl, repo)
	form := url.Values{}
	form.Set("date", "") // missing date
	form.Set("memo", "Test")
	form.Set("amount", "10")
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}
