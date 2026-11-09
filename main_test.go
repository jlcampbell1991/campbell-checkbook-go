package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	rootHandler(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	bodyStr := string(body)
	// Check navigation links
	expectedLinks := []string{"Transactions", "Pay Days", "Monthly Budget", "Inactive Line Items", "New Account"}
	for _, link := range expectedLinks {
		if !strings.Contains(bodyStr, link) {
			t.Errorf("response missing navigation link %s", link)
		}
	}
	// Check body content
	if !strings.Contains(bodyStr, "Welcome to Checkbook") {
		t.Errorf("response missing body content")
	}
}
