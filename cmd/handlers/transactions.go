package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

func IsWithdrawal(amount float64) bool    { return amount < 0 }
func CurrencyFromDouble(v float64) string { return fmt.Sprintf("%.2f", v) }

// transactionService holds dependencies for transaction handlers.
type transactionService struct {
	repo db.TransactionRepository
	tmpl *template.Template
}

// TransactionRoutes returns a router for transaction-related endpoints.
func TransactionRoutes(t *template.Template, repo db.TransactionRepository) chi.Router {
	s := &transactionService{repo: repo, tmpl: t}
	r := chi.NewRouter()
	r.Get("/", s.index)
	r.Get("/new", s.new)
	r.Post("/", s.create)
	// Additional routes (edit, delete, etc.) could be added here.
	return r
}

// index renders the transactions list page.
func (s *transactionService) index(w http.ResponseWriter, r *http.Request) {
	page := 0
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v >= 0 {
			page = v
		}
	}
	transactions, err := s.repo.Index(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nextURL := fmt.Sprintf("/transactions?page=%d", page+1)
	data := struct {
		Transactions []domain.Transaction
		NextURL      string
	}{
		Transactions: transactions,
		NextURL:      nextURL,
	}
	if err := s.tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// new renders the form to create a new transaction.
func (s *transactionService) new(w http.ResponseWriter, r *http.Request) {
	data := struct {
		AccountID  string
		LineItemID string
	}{
		AccountID:  r.URL.Query().Get("account_id"),
		LineItemID: r.URL.Query().Get("line_item_id"),
	}
	if err := s.tmpl.ExecuteTemplate(w, "transactions_add.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// create processes the submitted transaction form.
func (s *transactionService) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dateStr := r.FormValue("date")
	memo := r.FormValue("memo")
	amtStr := r.FormValue("amount")
	cat := r.FormValue("category")
	lineItemID := r.FormValue("line_item_id")
	accountID := r.FormValue("account_id")
	if dateStr == "" || memo == "" || amtStr == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}
	amt, err := strconv.ParseFloat(amtStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}
	tx := domain.Transaction{Date: d, Memo: memo, Amount: amt, Category: cat, LineItemID: lineItemID, AccountID: accountID}
	if _, err := s.repo.Create(r.Context(), tx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	redirect := "/"
	if lineItemID != "" {
		redirect = fmt.Sprintf("/line-items/%s", lineItemID)
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

// Existing TransactionsIndex kept for backward compatibility (may be removed).
func TransactionsIndex(w http.ResponseWriter, r *http.Request) {
	// Delegate to the service index using the default stub repository and template parsing.
	// This function is retained to avoid breaking existing routes that may still reference it.
	// It parses templates directly as before.
	page := 0
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v >= 0 {
			page = v
		}
	}
	transactions, err := db.TransactionRepo.Index(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nextURL := fmt.Sprintf("/transactions?page=%d", page+1)
	data := struct {
		Transactions []domain.Transaction
		NextURL      string
	}{
		Transactions: transactions,
		NextURL:      nextURL,
	}
	tmpl, err := template.ParseFiles(
		filepath.Join("views", "layouts", "application.tmpl"),
		filepath.Join("views", "layouts", "menu.tmpl"),
		filepath.Join("views", "transactions", "layout.tmpl"),
		filepath.Join("views", "transactions", "index.tmpl"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl = tmpl.Funcs(template.FuncMap{"IsWithdrawal": IsWithdrawal, "CurrencyFromDouble": CurrencyFromDouble})
	if err := tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
