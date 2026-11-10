package handlers

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "html/template"
    "io"
)

func TestTransactionsIndexTemplate(t *testing.T) {
    // Prepare sample data
    transactions = []Transaction{{ID: 1, AccountID: "A1", Description: "Test txn", Amount: 123.45}}
    // Ensure templates are loaded via init() in main package; tmpl is global variable in main package, not accessible here.
    // Instead, we create a new template parser for this test.
    // Load templates from the embedded FS used in main.
    // Reuse the same pattern.
    tmpl, err := template.ParseFiles("../../views/transactionsIndex.tmpl")
    if err != nil {
        t.Fatalf("failed to parse templates: %v", err)
    }
    handler := listTransactionsHandler(tmpl)
    req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
    w := httptest.NewRecorder()
    handler(w, req)
    res := w.Result()
    if res.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(res.Body)
        t.Fatalf("expected status 200, got %d, body: %s", res.StatusCode, string(bodyBytes))
    }
    var buf bytes.Buffer
    buf.ReadFrom(res.Body)
    body := buf.String()
    if !bytes.Contains([]byte(body), []byte("Test txn")) {
        t.Fatalf("response does not contain transaction description: %s", body)
    }
    if !bytes.Contains([]byte(body), []byte("123.45")) {
        t.Fatalf("response does not contain transaction amount: %s", body)
    }
}
