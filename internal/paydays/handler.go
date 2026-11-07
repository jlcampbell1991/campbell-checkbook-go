package paydays

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// ----- Request structs -----
type CreatePayDayRequest struct {
	Date     string `json:"date"` // ISO 8601 date string
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type UpdatePayDayRequest struct {
	Date     *string `json:"date,omitempty"`
	Amount   *int64  `json:"amount,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

type CreateDepositRequest struct {
	PayDayID  int64  `json:"payday_id"`
	AccountID int64  `json:"account_id"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
}

type UpdateDepositRequest struct {
	PayDayID  *int64  `json:"payday_id,omitempty"`
	AccountID *int64  `json:"account_id,omitempty"`
	Amount    *int64  `json:"amount,omitempty"`
	Currency  *string `json:"currency,omitempty"`
}

// ----- Validation helpers -----
func parseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("date is required")
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func validatePayDayInput(date time.Time, amount int64, currency string) error {
	if date.IsZero() {
		return errors.New("date must not be zero")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if currency == "" {
		return errors.New("currency must not be empty")
	}
	return nil
}

func validateDepositInput(paydayID, accountID, amount int64, currency string) error {
	if paydayID <= 0 {
		return errors.New("payday_id must be positive")
	}
	if accountID <= 0 {
		return errors.New("account_id must be positive")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if currency == "" {
		return errors.New("currency must not be empty")
	}
	return nil
}

// ----- PayDay handlers -----
func createPayDayHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePayDayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		date, err := parseDate(req.Date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validatePayDayInput(date, req.Amount, req.Currency); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pd := &PayDay{Date: date, Amount: req.Amount, Currency: req.Currency}
		if err := repo.CreatePayDay(r.Context(), pd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(pd)
	}
}

func getPayDayHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		pd, err := repo.GetPayDayByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if pd == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(pd)
	}
}

func listPayDaysHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := repo.ListPayDays(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(list)
	}
}

func updatePayDayHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		pd, err := repo.GetPayDayByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if pd == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		var req UpdatePayDayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Date != nil {
			d, err := parseDate(*req.Date)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			pd.Date = d
		}
		if req.Amount != nil {
			pd.Amount = *req.Amount
		}
		if req.Currency != nil {
			pd.Currency = *req.Currency
		}
		if err := validatePayDayInput(pd.Date, pd.Amount, pd.Currency); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := repo.UpdatePayDay(r.Context(), pd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(pd)
	}
}

func deletePayDayHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := repo.DeletePayDay(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ----- Deposit handlers -----
func createDepositHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateDepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if err := validateDepositInput(req.PayDayID, req.AccountID, req.Amount, req.Currency); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Ensure PayDay exists
		pd, err := repo.GetPayDayByID(r.Context(), req.PayDayID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if pd == nil {
			http.Error(w, "payday not found", http.StatusBadRequest)
			return
		}
		dep := &PayDayDeposit{PayDayID: req.PayDayID, AccountID: req.AccountID, Amount: req.Amount, Currency: req.Currency}
		if err := repo.CreateDeposit(r.Context(), dep); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(dep)
	}
}

func getDepositHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "deposit_id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		dep, err := repo.GetDepositByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if dep == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(dep)
	}
}

func listDepositsHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paydayIDStr := chi.URLParam(r, "id")
		paydayID, err := strconv.ParseInt(paydayIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid payday id", http.StatusBadRequest)
			return
		}
		list, err := repo.ListDepositsByPayDay(r.Context(), paydayID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(list)
	}
}

func updateDepositHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "deposit_id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		dep, err := repo.GetDepositByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if dep == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		var req UpdateDepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.PayDayID != nil {
			dep.PayDayID = *req.PayDayID
		}
		if req.AccountID != nil {
			dep.AccountID = *req.AccountID
		}
		if req.Amount != nil {
			dep.Amount = *req.Amount
		}
		if req.Currency != nil {
			dep.Currency = *req.Currency
		}
		if err := validateDepositInput(dep.PayDayID, dep.AccountID, dep.Amount, dep.Currency); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := repo.UpdateDeposit(r.Context(), dep); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(dep)
	}
}

func deleteDepositHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "deposit_id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := repo.DeleteDeposit(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
