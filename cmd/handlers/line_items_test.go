package handlers

import (
    "bytes"
    "context"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
    "text/template"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

type mockRepo struct{}

func (m *mockRepo) Show(ctx context.Context, id string) (domain.LineItem, error) {
    return domain.LineItem{ID: id, AccountID: "acc-1", Name: "Test LineItem", Balance: 123.45}, nil
}
func (m *mockRepo) GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error) {
    return []domain.Transaction{{ID: "tx1", LineItemID: lineItemID, Amount: 10.0, Category: "Food", Date: "2023-01-01"}}, nil
}
func (m *mockRepo) GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error) {
    return []domain.PayDayDeposit{{ID: "dep1", LineItemID: lineItemID, Amount: 200.0, Date: "2023-01-02"}}, nil
}

func TestShowHandler(t *testing.T) {
    // minimal template set containing the show template and layout placeholder
    tmplContent := `{{define "application.tmpl"}}{{template "content" .}}{{end}}` +
        `{{define "content"}}{{end}}` +
        `{{define "line_items_show.tmpl"}}LineItem: {{.LineItem.Name}} Next: {{.NextURL}}{{end}}`
    tmpl, err := template.New("test").Parse(tmplContent)
    if err != nil {
        t.Fatalf("parse template: %v", err)
    }
    repo := &mockRepo{}
    r := LineItemRoutes(tmpl, repo)

    req := httptest.NewRequest("GET", "/123?page=2", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp.StatusCode)
    }
    body, _ := ioutil.ReadAll(resp.Body)
    if !bytes.Contains(body, []byte("Test LineItem")) {
        t.Errorf("response missing line item name")
    }
    if !bytes.Contains(body, []byte("/line-items/123?page=3")) {
        t.Errorf("response missing next URL")
    }
}
