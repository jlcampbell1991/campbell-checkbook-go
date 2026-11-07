package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"campbell-checkbook/internal/accounts"
	"campbell-checkbook/internal/lineitems"
	"campbell-checkbook/internal/reconcile"
	"campbell-checkbook/internal/server"
	"campbell-checkbook/internal/transactions"
)

type mockAccountsRepo struct {
	data map[int64]*accounts.Account
	next int64
}

func newMockAccountsRepo() *mockAccountsRepo {
	return &mockAccountsRepo{data: make(map[int64]*accounts.Account), next: 1}
}

func (m *mockAccountsRepo) Create(ctx context.Context, acct *accounts.Account) error {
	acct.ID = m.next
	m.next++
	m.data[acct.ID] = acct
	return nil
}
func (m *mockAccountsRepo) GetByID(ctx context.Context, id int64) (*accounts.Account, error) {
	if a, ok := m.data[id]; ok {
		return a, nil
	}
	return nil, nil
}
func (m *mockAccountsRepo) Update(ctx context.Context, acct *accounts.Account) error { return nil }
func (m *mockAccountsRepo) Delete(ctx context.Context, id int64) error {
	delete(m.data, id)
	return nil
}
func (m *mockAccountsRepo) List(ctx context.Context) ([]*accounts.Account, error) {
	list := make([]*accounts.Account, 0, len(m.data))
	for _, a := range m.data {
		list = append(list, a)
	}
	return list, nil
}

type mockLineItemsRepo struct {
	data map[int64]*lineitems.LineItem
	next int64
}

func newMockLineItemsRepo() *mockLineItemsRepo {
	return &mockLineItemsRepo{data: make(map[int64]*lineitems.LineItem), next: 1}
}
func (m *mockLineItemsRepo) Create(ctx context.Context, li *lineitems.LineItem) error {
	li.ID = m.next
	m.next++
	m.data[li.ID] = li
	return nil
}
func (m *mockLineItemsRepo) GetByID(ctx context.Context, id int64) (*lineitems.LineItem, error) {
	if v, ok := m.data[id]; ok {
		return v, nil
	}
	return nil, nil
}
func (m *mockLineItemsRepo) Update(ctx context.Context, li *lineitems.LineItem) error { return nil }
func (m *mockLineItemsRepo) Delete(ctx context.Context, id int64) error {
	delete(m.data, id)
	return nil
}
func (m *mockLineItemsRepo) List(ctx context.Context) ([]*lineitems.LineItem, error) {
	list := make([]*lineitems.LineItem, 0, len(m.data))
	for _, v := range m.data {
		list = append(list, v)
	}
	return list, nil
}

type mockTransactionsRepo struct {
	imported int
}

func (m *mockTransactionsRepo) Create(ctx context.Context, tx *transactions.Transaction) error {
	return nil
}
func (m *mockTransactionsRepo) GetByID(ctx context.Context, id int64) (*transactions.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionsRepo) Update(ctx context.Context, tx *transactions.Transaction) error {
	return nil
}
func (m *mockTransactionsRepo) Delete(ctx context.Context, id int64) error { return nil }
func (m *mockTransactionsRepo) List(ctx context.Context) ([]*transactions.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionsRepo) ImportCSV(ctx context.Context, r io.Reader) (int, error) {
	data, _ := io.ReadAll(r)
	lines := bytes.Count(data, []byte("\n"))
	if lines > 0 {
		lines--
	}
	m.imported = lines
	return lines, nil
}

type mockReconcileRepo struct {
	data map[int64]*reconcile.ImportJob
	next int64
}

func newMockReconcileRepo() *mockReconcileRepo {
	return &mockReconcileRepo{data: make(map[int64]*reconcile.ImportJob), next: 1}
}
func (m *mockReconcileRepo) Create(ctx context.Context, job *reconcile.ImportJob) error {
	job.ID = m.next
	m.next++
	m.data[job.ID] = job
	return nil
}
func (m *mockReconcileRepo) GetByID(ctx context.Context, id int64) (*reconcile.ImportJob, error) {
	if j, ok := m.data[id]; ok {
		return j, nil
	}
	return nil, nil
}
func (m *mockReconcileRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	if j, ok := m.data[id]; ok {
		j.Status = status
	}
	return nil
}
func (m *mockReconcileRepo) List(ctx context.Context) ([]*reconcile.ImportJob, error) {
	list := make([]*reconcile.ImportJob, 0, len(m.data))
	for _, j := range m.data {
		list = append(list, j)
	}
	return list, nil
}

func TestFullStack(t *testing.T) {
	// Setup router with mock repos
	r := server.NewRouter()
	accountsRepo := newMockAccountsRepo()
	lineItemsRepo := newMockLineItemsRepo()
	txRepo := &mockTransactionsRepo{}
	reconcileRepo := newMockReconcileRepo()

	accounts.RegisterRoutes(r, accountsRepo)
	lineitems.RegisterRoutes(r, lineItemsRepo)
	transactions.RegisterRoutes(r, txRepo)
	reconcile.RegisterRoutes(r, reconcileRepo)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// 1. Create account
	acctPayload := map[string]interface{}{"name": "Test", "balance": 1000, "currency": "USD"}
	acctBody, _ := json.Marshal(acctPayload)
	resp, err := http.Post(ts.URL+"/accounts/", "application/json", bytes.NewReader(acctBody))
	if err != nil {
		t.Fatalf("POST /accounts error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var acctResp accounts.Account
	json.NewDecoder(resp.Body).Decode(&acctResp)
	resp.Body.Close()
	if acctResp.ID == 0 {
		t.Fatalf("expected account ID set")
	}

	// 2. Create line item linked to account
	liPayload := map[string]interface{}{
		"account_id":  acctResp.ID,
		"description": "Test item",
		"amount":      500,
		"currency":    "USD",
		"date":        time.Now().Format("2006-01-02"),
	}
	liBody, _ := json.Marshal(liPayload)
	resp, err = http.Post(ts.URL+"/line-items/", "application/json", bytes.NewReader(liBody))
	if err != nil {
		t.Fatalf("POST /line-items error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 line item, got %d", resp.StatusCode)
	}
	var liResp lineitems.LineItem
	json.NewDecoder(resp.Body).Decode(&liResp)
	resp.Body.Close()
	if liResp.ID == 0 {
		t.Fatalf("expected line item ID set")
	}

	// 3. Import transactions CSV
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "tx.csv")
	csv := "account_id,line_item_id,amount,currency,date,description\n" +
		fmt.Sprintf("%d,%d,500,USD,%s,desc\n", acctResp.ID, liResp.ID, time.Now().Format("2006-01-02"))
	part.Write([]byte(csv))
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/transactions/import", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST /transactions/import error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 import, got %d", resp.StatusCode)
	}
	var importResp struct {
		Imported int `json:"imported"`
	}
	json.NewDecoder(resp.Body).Decode(&importResp)
	resp.Body.Close()
	if importResp.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", importResp.Imported)
	}

	// 4. Create reconcile job
	jobPayload := map[string]string{"file_path": "dummy.csv"}
	jobBody, _ := json.Marshal(jobPayload)
	resp, err = http.Post(ts.URL+"/reconcile/jobs", "application/json", bytes.NewReader(jobBody))
	if err != nil {
		t.Fatalf("POST /reconcile/jobs error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 job, got %d", resp.StatusCode)
	}
	var jobResp reconcile.ImportJob
	json.NewDecoder(resp.Body).Decode(&jobResp)
	resp.Body.Close()
	if jobResp.ID == 0 {
		t.Fatalf("expected job ID set")
	}
}
