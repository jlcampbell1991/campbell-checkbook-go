package lineitems

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type memRepo struct {
	data map[int64]*LineItem
	next int64
}

func newMemRepo() *memRepo {
	return &memRepo{data: make(map[int64]*LineItem), next: 1}
}

func (m *memRepo) Create(ctx context.Context, li *LineItem) error {
	li.ID = m.next
	m.next++
	now := time.Now().UTC()
	li.CreatedAt = now
	li.UpdatedAt = now
	// copy to store
	copyLi := *li
	m.data[li.ID] = &copyLi
	return nil
}

func (m *memRepo) GetByID(ctx context.Context, id int64) (*LineItem, error) {
	if li, ok := m.data[id]; ok {
		return li, nil
	}
	return nil, nil
}

func (m *memRepo) Update(ctx context.Context, li *LineItem) error {
	if _, ok := m.data[li.ID]; !ok {
		return nil
	}
	li.UpdatedAt = time.Now().UTC()
	copyLi := *li
	m.data[li.ID] = &copyLi
	return nil
}

func (m *memRepo) Delete(ctx context.Context, id int64) error {
	delete(m.data, id)
	return nil
}

func (m *memRepo) List(ctx context.Context) ([]*LineItem, error) {
	var list []*LineItem
	for _, li := range m.data {
		list = append(list, li)
	}
	return list, nil
}

func TestLineItemHandlers(t *testing.T) {
	repo := newMemRepo()
	r := chi.NewRouter()
	RegisterRoutes(r, repo)
	server := httptest.NewServer(r)
	defer server.Close()

	// Create line item
	createPayload := CreateLineItemRequest{AccountID: 1, Description: "Test item", Amount: 1234, Currency: "USD", Date: "2023-01-02"}
	body, _ := json.Marshal(createPayload)
	resp, err := http.Post(server.URL+"/line-items/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var created LineItem
	_ = json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	if created.ID == 0 {
		t.Fatalf("expected non-zero ID")
	}

	// Get line item
	getResp, err := http.Get(server.URL + "/line-items/" + strconv.FormatInt(created.ID, 10))
	if err != nil {
		t.Fatalf("GET error: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", getResp.StatusCode)
	}
	var fetched LineItem
	_ = json.NewDecoder(getResp.Body).Decode(&fetched)
	getResp.Body.Close()
	if fetched.Description != "Test item" {
		t.Fatalf("expected description 'Test item', got %s", fetched.Description)
	}

	// Update line item
	newDesc := "Updated item"
	updPayload := UpdateLineItemRequest{Description: &newDesc}
	updBody, _ := json.Marshal(updPayload)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/line-items/"+strconv.FormatInt(created.ID, 10), bytes.NewReader(updBody))
	req.Header.Set("Content-Type", "application/json")
	updResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT error: %v", err)
	}
	if updResp.StatusCode != http.StatusOK {
		t.Fatalf("PUT expected 200, got %d", updResp.StatusCode)
	}
	var updated LineItem
	_ = json.NewDecoder(updResp.Body).Decode(&updated)
	updResp.Body.Close()
	if updated.Description != newDesc {
		t.Fatalf("expected updated description %s, got %s", newDesc, updated.Description)
	}

	// List line items
	listResp, err := http.Get(server.URL + "/line-items/")
	if err != nil {
		t.Fatalf("LIST error: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("LIST expected 200, got %d", listResp.StatusCode)
	}
	var list []LineItem
	_ = json.NewDecoder(listResp.Body).Decode(&list)
	listResp.Body.Close()
	if len(list) != 1 {
		t.Fatalf("expected list length 1, got %d", len(list))
	}

	// Delete line item
	delReq, _ := http.NewRequest(http.MethodDelete, server.URL+"/line-items/"+strconv.FormatInt(created.ID, 10), nil)
	delResp, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatalf("DELETE error: %v", err)
	}
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE expected 204, got %d", delResp.StatusCode)
	}
	delResp.Body.Close()
}
