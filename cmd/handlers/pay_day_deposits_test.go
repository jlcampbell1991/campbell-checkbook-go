package handlers

import (
    "context"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    "html/template"
)

type testPayDayRepo struct{}

func (t *testPayDayRepo) GetAll(ctx context.Context) ([]domain.PayDay, error) {
    return []domain.PayDay{{Id: domain.PayDayId{Value: "pd1"}, Name: "Test PayDay", Total: 1000, AccountId: "acc1"}}, nil
}
func (t *testPayDayRepo) GetAllInactive(ctx context.Context) ([]domain.PayDay, error) { return nil, nil }
func (t *testPayDayRepo) GetById(ctx context.Context, id string) (domain.PayDay, error) { return domain.PayDay{}, nil }
func (t *testPayDayRepo) Insert(ctx context.Context, req *domain.PayDayReq) (*domain.PayDay, error) { return nil, nil }
func (t *testPayDayRepo) Update(ctx context.Context, req *domain.PayDayReq, id string) (*domain.PayDay, error) { return nil, nil }
func (t *testPayDayRepo) Deactivate(ctx context.Context, id string) error { return nil }
func (t *testPayDayRepo) Reactivate(ctx context.Context, id string) error { return nil }

type testLineItemRepo struct{}
func (t *testLineItemRepo) Show(ctx context.Context, id string) (domain.LineItem, error) { return domain.LineItem{}, nil }
func (t *testLineItemRepo) GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error) { return nil, nil }
func (t *testLineItemRepo) GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error) { return nil, nil }
func (t *testLineItemRepo) Insert(ctx context.Context, li domain.LineItem) (string, error) { return "", nil }

func TestPayDayDepositNew(t *testing.T) {
    tmpl, _ := template.New("test").Parse("{{define \"pay_day_deposits_add.tmpl\"}}test{{end}}")
    r := PayDayDepositRoutes(tmpl, &testPayDayRepo{}, &db.StubPayDayDepositRepo{}, &testLineItemRepo{})
    req := httptest.NewRequest("GET", "/new?line-item-id=item1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "test") {
        t.Fatalf("response does not contain expected content")
    }
}

func TestPayDayDepositCreate(t *testing.T) {
    tmpl, _ := template.New("test").Parse("{{define \"pay_day_deposits_add.tmpl\"}}test{{end}}")
    // Use a stub deposit repo that records calls.
    depositRepo := &db.StubPayDayDepositRepo{}
    r := PayDayDepositRoutes(tmpl, &testPayDayRepo{}, depositRepo, &testLineItemRepo{})
    form := "amount=100&payDay=pd1&lineItem=item1"
    req := httptest.NewRequest("POST", "/", strings.NewReader(form))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusSeeOther {
        t.Fatalf("expected 303 SeeOther, got %d", w.Code)
    }
    expectedLocation := "/line-items/item1"
    if loc := w.Header().Get("Location"); loc != expectedLocation {
        t.Fatalf("expected redirect to %s, got %s", expectedLocation, loc)
    }
}
