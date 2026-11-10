package main

import (
    "bytes"
    "fmt"
    "testing"
)

func TestLineItemFormTemplate(t *testing.T) {
    // Ensure templates are loaded via init()
    data := struct {
        AccountID string
    }{AccountID: "acc-123"}
    var buf bytes.Buffer
    if err := tmpl.ExecuteTemplate(&buf, "lineItemForm", data); err != nil {
        t.Fatalf("lineItemForm template execution failed: %v", err)
    }
    // Verify that the form action contains the account ID
    expectedAction := fmt.Sprintf("/accounts/%s/line-items", data.AccountID)
    if !bytes.Contains(buf.Bytes(), []byte(expectedAction)) {
        t.Fatalf("output does not contain form action %s: %s", expectedAction, buf.String())
    }
    // Verify hidden input contains account ID
    hidden := fmt.Sprintf("value=\"%s\"", data.AccountID)
    if !bytes.Contains(buf.Bytes(), []byte(hidden)) {
        t.Fatalf("output does not contain hidden accountId input %s: %s", hidden, buf.String())
    }
}
