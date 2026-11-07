package paydays

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockRepo struct {
	created  []*PayDay
	deposits []*PayDayDeposit
	err      error
}

// PayDay methods
func (m *mockRepo) CreatePayDay(ctx context.Context, pd *PayDay) error {
	if m.err != nil {
		return m.err
	}
	// simulate DB assigning ID and timestamps
	pd.ID = int64(len(m.created) + 1)
	now := time.Now().UTC()
	pd.CreatedAt = now
	pd.UpdatedAt = now
	m.created = append(m.created, pd)
	return nil
}
func (m *mockRepo) GetPayDayByID(ctx context.Context, id int64) (*PayDay, error) {
	for _, pd := range m.created {
		if pd.ID == id {
			return pd, nil
		}
	}
	return nil, nil
}
func (m *mockRepo) UpdatePayDay(ctx context.Context, pd *PayDay) error { return nil }
func (m *mockRepo) DeletePayDay(ctx context.Context, id int64) error   { return nil }
func (m *mockRepo) ListPayDays(ctx context.Context) ([]*PayDay, error) { return m.created, nil }

// Deposit methods (no‑op for this test)
func (m *mockRepo) CreateDeposit(ctx context.Context, d *PayDayDeposit) error { return nil }
func (m *mockRepo) GetDepositByID(ctx context.Context, id int64) (*PayDayDeposit, error) {
	return nil, nil
}
func (m *mockRepo) UpdateDeposit(ctx context.Context, d *PayDayDeposit) error { return nil }
func (m *mockRepo) DeleteDeposit(ctx context.Context, id int64) error         { return nil }
func (m *mockRepo) ListDepositsByPayDay(ctx context.Context, paydayID int64) ([]*PayDayDeposit, error) {
	return nil, nil
}

func TestCreatePayDayHandler(t *testing.T) {
	repo := &mockRepo{}
	handler := createPayDayHandler(repo)

	payload := `{"date":"2023-12-01","amount":5000,"currency":"USD"}`
	req := httptest.NewRequest(http.MethodPost, "/paydays", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
	var resp PayDay
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID == 0 {
		t.Fatalf("expected non‑zero ID")
	}
	if resp.Amount != 5000 || resp.Currency != "USD" {
		t.Fatalf("unexpected response fields: %+v", resp)
	}
	if resp.Date.IsZero() {
		t.Fatalf("date not set")
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected repository to have 1 record, got %d", len(repo.created))
	}
}
