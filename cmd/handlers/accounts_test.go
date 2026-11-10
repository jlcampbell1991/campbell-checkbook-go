package handlers

import (
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "embed"
    "html/template"
    "github.com/go-chi/chi/v5"
)

//go:embed views/layouts/*.tmpl
var testTemplateFS embed.FS

func setupRouter(t *testing.T) (*chi.Mux, *template.Template) {
    tmpl, err := template.ParseFS(testTemplateFS, "views/layouts/*.tmpl")
    if err != nil {
        t.Fatalf("failed to parse test templates: %v", err)
    }
    r := chi.NewRouter()
    r.Get("/accounts/new", GetNewAccountForm(tmpl))
    r.Post("/accounts", PostAccount)
    return r, tmpl
}

func TestGetNewAccountForm(t *testing.T) {
    r, _ := setupRouter(t)
    req := httptest.NewRequest(http.MethodGet, "/accounts/new", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
    body := w.Body.String()
    if !strings.Contains(body, "<form") {
        t.Fatalf("response does not contain form: %s", body)
    }
}

func TestPostAccount(t *testing.T) {
    r, _ := setupRouter(t)
    form := strings.NewReader("name=TestAccount")
    req := httptest.NewRequest(http.MethodPost, "/accounts", form)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusSeeOther {
        t.Fatalf("expected status %d, got %d", http.StatusSeeOther, w.Code)
    }
    loc := w.Header().Get("Location")
    if loc != "/" {
        t.Fatalf("expected redirect to '/', got %s", loc)
    }
    // Ensure body is empty (optional)
    b, _ := io.ReadAll(w.Body)
    if len(b) != 0 {
        t.Fatalf("expected empty body on redirect, got %s", string(b))
    }
}
