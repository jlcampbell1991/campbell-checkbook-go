package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// PromptRequest represents the request payload sent to the Hermes LLM service.
type PromptRequest struct {
	// Prompt is the user supplied text that the LLM should complete.
	Prompt string `json:"prompt"`
	// Model optionally selects the LLM model to use. Empty means default.
	Model string `json:"model,omitempty"`
	// Temperature controls randomness. Omit for default.
	Temperature *float64 `json:"temperature,omitempty"`
	// MaxTokens limits the generated completion length.
	MaxTokens *int `json:"max_tokens,omitempty"`
}

// PromptResponse represents the JSON response returned by Hermes.
type PromptResponse struct {
	// Completion contains the generated text.
	Completion string `json:"completion"`
	// ID is an optional identifier for the request.
	ID string `json:"id,omitempty"`
	// Model echoes the model used.
	Model string `json:"model,omitempty"`
}

// HermesClient provides methods to interact with the Hermes HTTP API.
type HermesClient struct {
	baseURL    string
	httpClient *http.Client
	// maxRetries defines how many times a request will be retried on transient errors.
	maxRetries int
	// retryDelay is the base delay between retries.
	retryDelay time.Duration
}

// NewClient creates a new HermesClient. If httpClient is nil a default client with a timeout is used.
func NewClient(baseURL string, httpClient *http.Client) *HermesClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	// Ensure baseURL does not end with a slash to simplify endpoint concatenation.
	baseURL = strings.TrimRight(baseURL, "/")
	return &HermesClient{baseURL: baseURL, httpClient: httpClient, maxRetries: 2, retryDelay: 200 * time.Millisecond}
}

// Prompt sends a PromptRequest to the Hermes service and returns the parsed PromptResponse.
func (c *HermesClient) Prompt(ctx context.Context, req PromptRequest) (PromptResponse, error) {
	var result PromptResponse
	// Encode request payload.
	payload, err := json.Marshal(req)
	if err != nil {
		return result, fmt.Errorf("failed to marshal PromptRequest: %w", err)
	}

	endpoint := c.baseURL + "/prompt"
	// Simple retry loop.
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
		if err != nil {
			return result, fmt.Errorf("failed to create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			// Network error – retry if attempts left.
			if attempt < c.maxRetries {
				time.Sleep(c.retryDelay)
				continue
			}
			return result, fmt.Errorf("request error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// Read body for debugging.
			bodyBytes, _ := io.ReadAll(resp.Body)
			errMsg := fmt.Sprintf("unexpected status %s: %s", resp.Status, string(bodyBytes))
			// Retry on 5xx errors.
			if resp.StatusCode >= 500 && attempt < c.maxRetries {
				time.Sleep(c.retryDelay)
				continue
			}
			return result, fmt.Errorf(errMsg)
		}

		// Decode successful response.
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, fmt.Errorf("failed to decode response: %w", err)
		}
		return result, nil
	}
	// Should not reach here.
	return result, fmt.Errorf("exhausted retries for Prompt request")
}
