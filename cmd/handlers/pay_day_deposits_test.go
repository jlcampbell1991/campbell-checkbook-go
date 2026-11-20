//go:build ignore

package handlers

import (
    "context"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    "html/template"
    "os"
)

// stubPayDayRepo implements db.PayDayRepository for tests.
type stubPayDayRepo struct{ payDays []domain.PayDay }

func (s *stubPayDayRepo) GetAll(ctx context.Context) ([]domain.PayDay, error) { return s.payDays, nil }
func (s *stubPayDayRepo) GetAllInactive(ctx context.Context) ([]domain.PayDay, error) { return nil, nil }
func (s *stubPayDayRepo) GetById(ctx context.Context, id string) (domain.PayDay, error) { return domain.PayDay{}, nil }
func (s *stubPayDayRepo) Insert(ctx context.Context, req *domain.PayDayReq) (*domain.PayDay, error) { return nil, nil }
func (s *stubPayDayRepo) Update(ctx context.Context, req *domain.PayDayReq, id string) (*domain.PayDay, error) { return nil, nil }
func (s *stubPayDayRepo) Deactivate(ctx context.Context, id string) error { return nil }
func (s *stubPayDayRepo) Reactivate(ctx context.Context, id string) error { return nil }

// stubDepositRepo records Upsert calls.
type stubDepositRepo struct{ called bool; dep domain.PayDayDeposit; pdID string }

func (s *stubDepositRepo) GetAllByPayDay(ctx context.Context, payDayId domain.PayDayId) ([]domain.PayDayDeposit, error) { return nil, nil }
func (s *stubDepositRepo) Upsert(ctx context.Context, dep domain.PayDayDeposit, payDayId domain.PayDayId) error {
    s.called = true
    s.dep = dep
    s.pdID = payDayId.Value
    return nil
}
func (s *stubDepositRepo) GetById(ctx context.Context, id string) (domain.PayDayDeposit, error) { return domain.PayDayDeposit{}, nil }
func (s *stubDepositRepo) Destroy(ctx context.Context, id string) error { return nil }

// stubLineItemRepo satisfies interface but not used.
type stubLineItemRepo struct{}
func (s *stubLineItemRepo) Insert(ctx context.Context, li domain.LineItem) (string, error) { return "", nil }
func (s *stubLineItemRepo) Show(ctx context.Context, id string) (domain.LineItem, error) { return domain.LineItem{}, nil }
func (s *stubLineItemRepo) GetTransactions(ctx context.Context, lineItemID string, page int) ([]domain.Transaction, error) { return nil, nil }
func (s *stubLineItemRepo) GetPayDayDeposits(ctx context.Context, lineItemID string) ([]domain.PayDayDeposit, error) { return nil, nil }


func TestPayDayDepositNewHandler(t *testing.T) {
    tmpl := template.Must(template.ParseFS(os.DirFS("../.."), "views/layouts/*.tmpl", "views/pay_day_deposits/*.tmpl"))
    pdRepo := &stubPayDayRepo{payDays: []domain.PayDay{{Id: domain.PayDayId{Value: "pd1"}, Name: "Test PayDay"}}}
    depRepo := &stubDepositRepo{}
    liRepo := &stubLineItemRepo{}
    handler := PayDayDepositRoutes(tmpl, pdRepo, depRepo, liRepo)

    req := httptest.NewRequest(http.MethodGet, "/new?line-item-id=li1", nil)
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "<form") {
        t.Fatalf("response does not contain form")
    }
    if !strings.Contains(w.Body.String(), "Test PayDay") {
        t.Fatalf("response does not contain pay day option")
    }
}

func TestPayDayDepositCreateHandler(t *testing.T) {
    tmpl := template.Must(template.ParseFS(os.DirFS("../.."), "views/layouts/*.tmpl", "views/pay_day_deposits/*.tmpl"))
    pdRepo := &stubPayDayRepo{}
    depRepo := &stubDepositRepo{}
    liRepo := &stubLineItemRepo{}
    handler := PayDayDepositRoutes(tmpl, pdRepo, depRepo, liRepo)

    form := "amount=100&payDay=pd1&lineItem=li1&date=2026-11-20"
    req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    if w.Code != http.StatusSeeOther {
        t.Fatalf("expected redirect status 303, got %d", w.Code)
    }
    if !depRepo.called {
        t.Fatalf("deposit repo Upsert not called")
    }
    if depRepo.pdID != "pd1" {
        t.Fatalf("unexpected payDayID %s", depRepo.pdID)
    }
    if depRepo.dep.LineItemID != "li1" || depRepo.dep.Amount != 100 {
        t.Fatalf("deposit fields not set correctly: %+v", depRepo.dep)
    }
}
