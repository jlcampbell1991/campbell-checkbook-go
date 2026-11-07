package accounts

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type UpdateAccountRequest struct {
	Name     *string `json:"name,omitempty"`
	Balance  *int64  `json:"balance,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

func validateAccountInput(name string, balance int64, currency string) error {
	if name == "" {
		return errors.New("name must not be empty")
	}
	if balance < 0 {
		return errors.New("balance must be non‑negative")
	}
	if currency == "" {
		return errors.New("currency must not be empty")
	}
	return nil
}

func createAccountHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if err := validateAccountInput(req.Name, req.Balance, req.Currency); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		acct := &Account{Name: req.Name, Balance: req.Balance, Currency: req.Currency}
		if err := repo.Create(r.Context(), acct); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(acct)
	}
}

func getAccountHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		acct, err := repo.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if acct == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(acct)
	}
}

func listAccountsHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := repo.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(accounts)
	}
}

func updateAccountHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		var req UpdateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		acct, err := repo.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if acct == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if req.Name != nil {
			acct.Name = *req.Name
		}
		if req.Balance != nil {
			acct.Balance = *req.Balance
		}
		if req.Currency != nil {
			acct.Currency = *req.Currency
		}
		if err := validateAccountInput(acct.Name, acct.Balance, acct.Currency); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := repo.Update(r.Context(), acct); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(acct)
	}
}

func deleteAccountHandler(repo Repository) http.HandlerFunc {
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
