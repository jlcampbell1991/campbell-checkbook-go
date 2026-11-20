package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "html/template"
    "github.com/go-chi/chi/v5"
    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

func TestPayDayNewHandler(t *testing.T) {
    // Parse minimal templates needed for the handler.
    tmpl, err := template.ParseGlob("../../views/layouts/*.tmpl")
    if err != nil {
        t.Fatalf("layout parse error: %v", err)
    }
    // Parse pay_day templates
    _, err = tmpl.ParseGlob("../../views/pay_days/*.tmpl")
    if err != nil {
        t.Fatalf("pay_days template parse error: %v", err)
    }
    if err != nil {
        t.Fatalf("template parse error: %v", err)
    }
    // Stub repos
    pdRepo := db.NewStubPayDayRepo()
    depRepo := db.NewStubPayDayDepositRepo()
    lineRepo := db.NewStubLineItemRepo()
    acctRepo := db.NewStubAccountRepo()

    r := chi.NewRouter()
    r.Mount("/pay-days", PayDayRoutes(tmpl, pdRepo, depRepo, lineRepo, acctRepo))

    req := httptest.NewRequest(http.MethodGet, "/pay-days/new", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}
