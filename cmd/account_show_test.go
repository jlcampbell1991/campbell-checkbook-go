package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAccountShowTemplate(t *testing.T) {
	// Ensure templates are loaded via init()
	data := struct {
		Account struct {
			ID   string
			Name string
		}
		LineItems []struct {
			Name    string
			Balance float64
		}
	}{
		Account: struct {
			ID   string
			Name string
		}{ID: "acc-123", Name: "Test Account"},
		LineItems: []struct {
			Name    string
			Balance float64
		}{{Name: "Item1", Balance: 10.5}},
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "accountShow", data); err != nil {
		t.Fatalf("accountShow template execution failed: %v", err)
	}
	// Verify that the account name appears
	if !bytes.Contains(buf.Bytes(), []byte(data.Account.Name)) {
		t.Fatalf("output does not contain account name: %s", buf.String())
	}
	// Verify that the delete link/button is present
	expected := fmt.Sprintf("/accounts/%s/delete", data.Account.ID)
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Fatalf("output does not contain delete link %s: %s", expected, buf.String())
	}
}
