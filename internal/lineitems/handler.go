package lineitems

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type CreateLineItemRequest struct {
	AccountID   int64  `json:"account_id"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Date        string `json:"date"` // ISO 8601 date string, e.g., "2023-01-02"
}

type UpdateLineItemRequest struct {
	AccountID   *int64  `json:"account_id,omitempty"`
	Description *string `json:"description,omitempty"`
	Amount      *int64  `json:"amount,omitempty"`
	Currency    *string `json:"currency,omitempty"`
	Date        *string `json:"date,omitempty"`
}

func validateLineItemInput(accountID int64, description string, amount int64, currency string, date time.Time) error {
	if accountID <= 0 {
		return errors.New("account_id must be positive")
	}
	if description == "" {
		return errors.New("description must not be empty")
	}
	if amount < 0 {
		return errors.New("amount must be non‑negative")
	}
	if currency == "" {
		return errors.New("currency must not be empty")
	}
	// date can be zero value? require non-zero
	if date.IsZero() {
		return errors.New("date must be provided")
	}
	return nil
}

func parseDateString(s string) (time.Time, error) {
	// Accept RFC3339 or simple date format
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// fallback to YYYY-MM-DD
	return time.Parse("2006-01-02", s)
}

func createLineItemHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateLineItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		dt, err := parseDateString(req.Date)
		if err != nil {
			http.Error(w, "invalid date format", http.StatusBadRequest)
			return
		}
		if err := validateLineItemInput(req.AccountID, req.Description, req.Amount, req.Currency, dt); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		li := &LineItem{AccountID: req.AccountID, Description: req.Description, Amount: req.Amount, Currency: req.Currency, Date: dt}
		if err := repo.Create(r.Context(), li); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(li)
	}
}

func getLineItemHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		li, err := repo.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if li == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(li)
	}
}

func listLineItemsHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := repo.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(items)
	}
}

func updateLineItemHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		var req UpdateLineItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		li, err := repo.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if li == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if req.AccountID != nil {
			li.AccountID = *req.AccountID
		}
		if req.Description != nil {
			li.Description = *req.Description
		}
		if req.Amount != nil {
			li.Amount = *req.Amount
		}
		if req.Currency != nil {
			li.Currency = *req.Currency
		}
		if req.Date != nil {
			dt, err := parseDateString(*req.Date)
			if err != nil {
				http.Error(w, "invalid date format", http.StatusBadRequest)
				return
			}
			li.Date = dt
		}
		if err := validateLineItemInput(li.AccountID, li.Description, li.Amount, li.Currency, li.Date); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := repo.Update(r.Context(), li); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(li)
	}
}

func deleteLineItemHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := repo.Delete(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
