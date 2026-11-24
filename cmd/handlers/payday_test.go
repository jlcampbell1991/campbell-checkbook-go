package handlers

import (
    "context"
    "bytes"
    "html/template"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// mockPayDayRepo provides a stub PayDayRepository for testing.
type mockPayDayRepo struct{}

func (m *mockPayDayRepo) GetAll(ctx context.Context) ([]domain.PayDay, error) { return []domain.PayDay{}, nil }
func (m *mockPayDayRepo) GetAllInactive(ctx context.Context) ([]domain.PayDay, error) { return []domain.PayDay{}, nil }
func (m *mockPayDayRepo) GetById(ctx context.Context, id string) (domain.PayDay, error) {
    return domain.PayDay{Id: domain.PayDayId{Value: id}, Name: "Test PayDay", Total: 100.0, AccountId: "acc-1"}, nil
}
func (m *mockPayDayRepo) Insert(ctx context.Context, req *domain.PayDayReq) (*domain.PayDay, error) {
    return &domain.PayDay{Id: domain.PayDayId{Value: "new-id"}, Name: req.Name, Total: req.Total, AccountId: req.AccountId}, nil
}
func (m *mockPayDayRepo) Update(ctx context.Context, req *domain.PayDayReq, id string) (*domain.PayDay, error) {
    return &domain.PayDay{Id: domain.PayDayId{Value: id}, Name: req.Name, Total: req.Total, AccountId: req.AccountId}, nil
}
func (m *mockPayDayRepo) Deactivate(ctx context.Context, id string) error { return nil }
func (m *mockPayDayRepo) Reactivate(ctx context.Context, id string) error { return nil }

type mockDepositRepo struct{}
func (m *mockDepositRepo) GetAllByPayDay(ctx context.Context, payDayId domain.PayDayId) ([]domain.PayDayDeposit, error) { return []domain.PayDayDeposit{}, nil }
func (m *mockDepositRepo) Upsert(ctx context.Context, dep domain.PayDayDeposit, payDayId domain.PayDayId) error { return nil }
func (m *mockDepositRepo) GetById(ctx context.Context, id string) (domain.PayDayDeposit, error) { return domain.PayDayDeposit{}, nil }
func (m *mockDepositRepo) Destroy(ctx context.Context, id string) error { return nil }

type mockLineItemRepo struct{}
func (m *mockLineItemRepo) Show(ctx context.Context, id string) (domain.LineItem, error) { return domain.LineItem{}, nil }
func (m *mockLineItemRepo) GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error) {
    return []domain.Transaction{}, nil
}
func (m *mockLineItemRepo) GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error) {
    return []domain.PayDayDeposit{}, nil
}
func (m *mockLineItemRepo) Insert(ctx context.Context, li domain.LineItem) (string, error) { return "stub-id", nil }

type mockAccountRepo2 struct{}

func (m *mockAccountRepo2) Insert(ctx context.Context, a domain.Account) (string, error) {
    return "acc-mock", nil
}

func TestNewPayDayForm(t *testing.T) {
    tmplContent := `{{define "application.tmpl"}}{{template "content" .}}{{end}}` +
        `{{define "content"}}{{end}}` +
        `{{define "pay_days_form.tmpl"}}<h2>{{.Title}}</h2>{{end}}`
    tmpl, err := template.New("test").Parse(tmplContent)
    if err != nil {
        t.Fatalf("parse template: %v", err)
    }
    r := PayDayRoutes(tmpl, &mockPayDayRepo{}, &mockDepositRepo{}, &mockLineItemRepo{}, &mockAccountRepo2{})
    req := httptest.NewRequest("GET", "/new", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp.StatusCode)
    }
    body, _ := ioutil.ReadAll(resp.Body)
    if !bytes.Contains(body, []byte("New Pay Day")) {
        t.Errorf("response does not contain expected title")
    }
}
