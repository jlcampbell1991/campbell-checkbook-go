package main

import (
    "bytes"
    "fmt"
    "testing"
)

func TestTransactionFormTemplate(t *testing.T) {
    // Ensure templates are loaded via init()
    var buf bytes.Buffer
    if err := tmpl.ExecuteTemplate(&buf, "transactionForm", nil); err != nil {
        t.Fatalf("transactionForm template execution failed: %v", err)
    }
    // Verify that the form action contains the expected path
    expectedAction := "/transactions"
    if !bytes.Contains(buf.Bytes(), []byte("action=\""+expectedAction+"\"")) {
        t.Fatalf("output does not contain form action %s: %s", expectedAction, buf.String())
    }
    // Verify that the form contains a field for description
    if !bytes.Contains(buf.Bytes(), []byte("name=\"description\"")) {
        t.Fatalf("output does not contain description input: %s", buf.String())
    }
    // Verify that the form contains a field for amount
    if !bytes.Contains(buf.Bytes(), []byte("name=\"amount\"")) {
        t.Fatalf("output does not contain amount input: %s", buf.String())
    }
    // Verify that the form contains a field for accountId
    if !bytes.Contains(buf.Bytes(), []byte("name=\"accountId\"")) {
        t.Fatalf("output does not contain accountId input: %s", buf.String())
    }
    // Optional: ensure hidden not required, just check presence
    fmt.Println(buf.String())
}
