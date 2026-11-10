package main

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

func TestLineItemAddTemplate(t *testing.T) {
	// Ensure templates are loaded via init()
	accountID := "test-account"
	// Build simple form HTML manually to simulate what handler would generate
	formHTML := fmt.Sprintf(`<form action="/accounts/%s/line-items" method="post"></form>`, accountID)
	page := pageData{Title: "Test", Body: template.HTML(formHTML)}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "lineItemAdd", page); err != nil {
		t.Fatalf("lineItemAdd template execution failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(fmt.Sprintf("/accounts/%s/line-items", accountID))) {
		t.Fatalf("output does not contain expected form action URL: %s", buf.String())
	}
}
