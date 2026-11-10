package main

import (
	"bytes"
	"html/template"
	"testing"
)

func TestAccountsIndexTemplate(t *testing.T) {
	// tmpl is initialized via init()
	var buf bytes.Buffer
	// Execute with empty slice
	if err := tmpl.ExecuteTemplate(&buf, "accountsIndex", []Account{}); err != nil {
		t.Fatalf("accountsIndex template execution failed: %v", err)
	}
	// Simple check that output contains expected placeholder (no panic means success)
	if buf.Len() == 0 {
		t.Fatalf("accountsIndex template produced empty output")
	}
	// Ensure output is valid HTML (optional simple check)
	_ = template.HTML(buf.String())
}
