//go:build ignore
// placeholder
package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestTransactionsIndexTemplate(t *testing.T) {
    // Prepare sample data
    transactions = []Transaction{{ID: 1, AccountID: "A1", Description: "Test txn", Amount: 123.45}}
    // Ensure templates are loaded via init()
    handler := listTransactionsHandler(tmpl)
    req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
    w := httptest.NewRecorder()
    handler(w, req)
    res := w.Result()
    if res.StatusCode != http.StatusOK {
        t.Fatalf("expected status 200, got %d", res.StatusCode)
    }
    var buf bytes.Buffer
    buf.ReadFrom(res.Body)
    body := buf.String()
    // Verify that the table contains the transaction data
    if !bytes.Contains([]byte(body), []byte("Test txn")) {
        t.Fatalf("response does not contain transaction description: %s", body)
    }
    if !bytes.Contains([]byte(body), []byte("123.45")) {
        t.Fatalf("response does not contain transaction amount: %s", body)
    }
}
