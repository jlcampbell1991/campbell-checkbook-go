package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestHermesClient_Prompt(t *testing.T) {
	// Expected request payload.
	expectedReq := PromptRequest{Prompt: "Hello, world!", Model: "test-model"}

	// Mock server that validates request and returns a response.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/prompt" {
			t.Fatalf("expected path /prompt, got %s", r.URL.Path)
		}
		var got PromptRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(got, expectedReq) {
			t.Fatalf("request payload mismatch.\nexpected: %+v\nactual:   %+v", expectedReq, got)
		}
		// Respond with JSON.
		resp := PromptResponse{Completion: "Hi there!", ID: "123", Model: "test-model"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, &http.Client{Timeout: 5 * time.Second})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	gotResp, err := client.Prompt(ctx, expectedReq)
	if err != nil {
		t.Fatalf("Prompt returned error: %v", err)
	}

	expectedResp := PromptResponse{Completion: "Hi there!", ID: "123", Model: "test-model"}
	if !reflect.DeepEqual(gotResp, expectedResp) {
		t.Fatalf("response mismatch.\nexpected: %+v\nactual:   %+v", expectedResp, gotResp)
	}
}
