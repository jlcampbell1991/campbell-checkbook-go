package handlers

import (
    "html/template"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/go-chi/chi/v5"
)


func TestRegisterRoutes(t *testing.T) {
    tmpl, err := template.ParseGlob("../../views/**/*.tmpl")
    if err != nil {
        t.Fatalf("template parse error: %v", err)
    }
    r := chi.NewRouter()
    RegisterRoutes(r, tmpl)

    // Test GET add form
    req := httptest.NewRequest(http.MethodGet, "/pay-day-deposits/add", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
    body := w.Body.String()
    if body == "" || !strings.Contains(body, "<form") {
        t.Fatalf("response does not contain form")
    }
}
