package handlers

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

type mockRepo struct{}

func (m *mockRepo) Index(ctx context.Context) ([]domain.Account, error) {
	bal := 100.0
	bud := 200.0
	acct := domain.Account{
		ID:      "1",
		Name:    "Test Account",
		Balance: &bal,
		Budget:  &bud,
		LineItems: []domain.LineItem{{
			ID:        "li1",
			Name:      "Item 1",
			Balance:   &bal,
			Budget:    &bud,
			Active:    true,
			AccountID: "1",
		}},
	}
	return []domain.Account{acct}, nil
}

func TestAccountsIndexHandler(t *testing.T) {
	// Parse templates with helpers
	funcMap := template.FuncMap{"Format": domain.Format, "IsWithdrawal": domain.IsWithdrawal}
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("../../views/**/*.tmpl")
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	handler := NewAccountsHandler(&mockRepo{}, tmpl)
	r := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// Directly call handler index via router
	handler.Router().ServeHTTP(r, req)

	if r.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", r.Code)
	}
	body := r.Body.String()
	if !bytes.Contains([]byte(body), []byte("Test Account")) {
		t.Fatalf("response does not contain account name: %s", body)
	}
	if !bytes.Contains([]byte(body), []byte("Item 1")) {
		t.Fatalf("response does not contain line item name: %s", body)
	}
}
