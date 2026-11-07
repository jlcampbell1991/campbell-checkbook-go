package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type memRepo struct {
	data map[int64]*Account
	next int64
}

func newMemRepo() *memRepo {
	return &memRepo{data: make(map[int64]*Account), next: 1}
}

func (m *memRepo) Create(ctx context.Context, acct *Account) error {
	acct.ID = m.next
	m.next++
	now := time.Now().UTC()
	acct.CreatedAt = now
	acct.UpdatedAt = now
	m.data[acct.ID] = acct
	return nil
}
func (m *memRepo) GetByID(ctx context.Context, id int64) (*Account, error) {
	if a, ok := m.data[id]; ok {
		return a, nil
	}
	return nil, nil
}
func (m *memRepo) Update(ctx context.Context, acct *Account) error {
	if _, ok := m.data[acct.ID]; !ok {
		return nil
	}
	acct.UpdatedAt = time.Now().UTC()
	m.data[acct.ID] = acct
	return nil
}
func (m *memRepo) Delete(ctx context.Context, id int64) error {
	delete(m.data, id)
	return nil
}
func (m *memRepo) List(ctx context.Context) ([]*Account, error) {
	var list []*Account
	for _, a := range m.data {
		list = append(list, a)
	}
	return list, nil
}

func TestAccountHandlers(t *testing.T) {
	repo := newMemRepo()
	r := chi.NewRouter()
	RegisterRoutes(r, repo)
	server := httptest.NewServer(r)
	defer server.Close()

	// Create account
	createPayload := CreateAccountRequest{Name: "Bob", Balance: 5000, Currency: "USD"}
	body, _ := json.Marshal(createPayload)
	resp, err := http.Post(server.URL+"/accounts/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var created Account
	_ = json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	if created.ID == 0 {
		t.Fatalf("expected non-zero ID")
	}

	// Get account
	getResp, err := http.Get(server.URL + "/accounts/" + strconv.FormatInt(created.ID, 10))
	if err != nil {
		t.Fatalf("GET error: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", getResp.StatusCode)
	}
	var fetched Account
	_ = json.NewDecoder(getResp.Body).Decode(&fetched)
	getResp.Body.Close()
	if fetched.Name != "Bob" {
		t.Fatalf("expected name Bob, got %s", fetched.Name)
	}

	// Update account
	newName := "Bob Updated"
	updPayload := UpdateAccountRequest{Name: &newName}
	updBody, _ := json.Marshal(updPayload)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/accounts/"+strconv.FormatInt(created.ID, 10), bytes.NewReader(updBody))
	req.Header.Set("Content-Type", "application/json")
	updResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT error: %v", err)
	}
	if updResp.StatusCode != http.StatusOK {
		t.Fatalf("PUT expected 200, got %d", updResp.StatusCode)
	}
	var updated Account
	_ = json.NewDecoder(updResp.Body).Decode(&updated)
	updResp.Body.Close()
	if updated.Name != newName {
		t.Fatalf("expected updated name %s, got %s", newName, updated.Name)
	}

	// List accounts
	listResp, err := http.Get(server.URL + "/accounts/")
	if err != nil {
		t.Fatalf("LIST error: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("LIST expected 200, got %d", listResp.StatusCode)
	}
	var list []Account
	_ = json.NewDecoder(listResp.Body).Decode(&list)
	listResp.Body.Close()
	if len(list) != 1 {
		t.Fatalf("expected list length 1, got %d", len(list))
	}

	// Delete account
	delReq, _ := http.NewRequest(http.MethodDelete, server.URL+"/accounts/"+strconv.FormatInt(created.ID, 10), nil)
	delResp, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatalf("DELETE error: %v", err)
	}
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE expected 204, got %d", delResp.StatusCode)
	}
	delResp.Body.Close()
}
