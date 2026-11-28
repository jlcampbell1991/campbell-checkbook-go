package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "html/template"
    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

func TestPayDayIndex(t *testing.T) {
    tmpl := template.New("").Parse(`{{define "application.tmpl"}}{{block "content" .}}{{end}}{{end}}`)
    // Parse the pay_days index template as well
    tmpl, _ = tmpl.ParseFiles("views/pay_days/index.tmpl")
    pdRepo := db.NewStubPayDayRepo()
    depRepo := db.NewStubPayDayDepositRepo()
    liRepo := db.LineItemRepo // existing stub
    acctRepo := db.NewStubAccountRepo()
    r := PayDayRoutes(tmpl, pdRepo, depRepo, liRepo, acctRepo)
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}
