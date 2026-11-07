package sessions

import (
	"bytes"
	"campbell-checkbook/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateAndValidateToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "testsecret"}
	userID := "42"
	token, err := GenerateToken(cfg, userID)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	gotID, err := ValidateToken(cfg, token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}
	if gotID != userID {
		t.Fatalf("expected userID %s, got %s", userID, gotID)
	}
}

func TestAuthMiddleware(t *testing.T) {
	cfg := &config.Config{JWTSecret: "secret"}
	// Create a valid token
	token, _ := GenerateToken(cfg, "99")

	// Protected handler returns 200
	protected := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := AuthMiddleware(cfg)(protected)

	// Test missing token
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing token, got %d", w.Code)
	}

	// Test invalid token
	req = httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", w.Code)
	}

	// Test valid token
	req = httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid token, got %d", w.Code)
	}
}

func TestLoginHandler(t *testing.T) {
	cfg := &config.Config{JWTSecret: "loginsecret"}
	// Prepare request body
	body, _ := json.Marshal(map[string]string{"username": "testuser", "password": "password"})
	req := httptest.NewRequest("POST", "/sessions/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler := loginHandler(cfg)
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp loginResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Token == "" {
		t.Fatalf("expected token in response")
	}
	// Validate token
	uid, err := ValidateToken(cfg, resp.Token)
	if err != nil {
		t.Fatalf("token validation failed: %v", err)
	}
	if uid != "1" { // testuser has UserID "1"
		t.Fatalf("expected userID 1, got %s", uid)
	}
}
