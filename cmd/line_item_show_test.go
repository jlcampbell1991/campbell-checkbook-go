package main

import (
    "bytes"
    "fmt"
    "testing"
)

func TestLineItemShowTemplate(t *testing.T) {
    // Ensure templates are loaded via init()
    data := struct {
        LineItem struct {
            ID      string
            Name    string
            Balance float64
        }
    }{}
    data.LineItem.ID = "li-123"
    data.LineItem.Name = "Test Item"
    data.LineItem.Balance = 99.99

    var buf bytes.Buffer
    if err := tmpl.ExecuteTemplate(&buf, "lineItemShow", data); err != nil {
        t.Fatalf("lineItemShow template execution failed: %v", err)
    }
    // Verify ID appears
    if !bytes.Contains(buf.Bytes(), []byte(data.LineItem.ID)) {
        t.Fatalf("output does not contain line item ID: %s", buf.String())
    }
    // Verify Name appears
    if !bytes.Contains(buf.Bytes(), []byte(data.LineItem.Name)) {
        t.Fatalf("output does not contain line item name: %s", buf.String())
    }
    // Verify Balance appears formatted
    expected := fmt.Sprintf("$%.2f", data.LineItem.Balance)
    if !bytes.Contains(buf.Bytes(), []byte(expected)) {
        t.Fatalf("output does not contain balance %s: %s", expected, buf.String())
    }
}
