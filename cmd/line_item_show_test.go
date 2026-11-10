package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLineItemShowTemplate(t *testing.T) {
	// Ensure templates are loaded via init()
	data := struct {
		Item struct {
			ID      string
			Name    string
			Balance float64
		}
		DeleteURL string
	}{
		Item: struct {
			ID      string
			Name    string
			Balance float64
		}{ID: "li-123", Name: "Test LineItem", Balance: 42.0},
		DeleteURL: "/line-items/li-123/delete",
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "lineItemShow", data); err != nil {
		t.Fatalf("lineItemShow template execution failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(data.Item.Name)) {
		t.Fatalf("output does not contain line item name: %s", buf.String())
	}
	expected := fmt.Sprintf("/line-items/%s/delete", data.Item.ID)
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Fatalf("output does not contain delete link %s: %s", expected, buf.String())
	}
}
