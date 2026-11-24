package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestTransactionNew(t *testing.T) {
	// Prepare template with required files (empty layout placeholder)
	tmpl, err := template.New("test").Parse(`{{define "application.tmpl"}}{{end}}`)
	if err != nil {
		t.Fatalf("template parse: %v", err)
	}
	// Add the add template to the set
	tmpl, err = tmpl.ParseFiles("../../views/transactions/add.tmpl")
	if err != nil {
		t.Fatalf("parse add tmpl: %v", err)
	}
	repo := &db.StubTransactionRepository{}
	r := chi.NewRouter()
	r.Mount("/transactions", TransactionRoutes(tmpl, repo))

	req := httptest.NewRequest(http.MethodGet, "/transactions/new?account_id=a1&line_item_id=l1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestTransactionCreate(t *testing.T) {
	tmpl, err := template.New("test").Parse(`{{define "application.tmpl"}}{{end}}`)
	if err != nil {
		t.Fatalf("template parse: %v", err)
	}
	tmpl, err = tmpl.ParseFiles("../../views/transactions/add.tmpl")
	if err != nil {
		t.Fatalf("parse add tmpl: %v", err)
	}
	repo := &db.StubTransactionRepository{}
	r := chi.NewRouter()
	r.Mount("/transactions", TransactionRoutes(tmpl, repo))

	form := url.Values{}
	form.Set("date", "2023-01-01")
	form.Set("memo", "Test transaction")
	form.Set("amount", "12.34")
	form.Set("category", "Misc")
	form.Set("account_id", "a1")
	form.Set("line_item_id", "l1")
	body := strings.NewReader(form.Encode())
	req := httptest.NewRequest(http.MethodPost, "/transactions", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect (303), got %d", w.Code)
	}
	// Expected redirect location
	expected := "/line-items/l1"
	if loc := w.Header().Get("Location"); loc != expected {
		t.Fatalf("expected redirect to %s, got %s", expected, loc)
	}
}

func TestTransactionCreateMissingFields(t *testing.T) {
	tmpl, _ := template.New("test").Parse(`{{define "application.tmpl"}}{{end}}`)
	repo := &db.StubTransactionRepository{}
	r := chi.NewRouter()
	r.Mount("/transactions", TransactionRoutes(tmpl, repo))

	form := url.Values{}
	form.Set("date", "") // missing
	form.Set("memo", "Test")
	form.Set("amount", "10")
	body := strings.NewReader(form.Encode())
	req := httptest.NewRequest(http.MethodPost, "/transactions", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
