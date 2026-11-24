package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// transactionService provides HTTP handlers for transaction operations.
type transactionService struct {
	repo db.TransactionRepository
	tmpl *template.Template
}

// TransactionRoutes returns a chi.Router with transaction‑related endpoints.
func TransactionRoutes(t *template.Template, repo db.TransactionRepository) chi.Router {
	s := &transactionService{repo: repo, tmpl: t}
	r := chi.NewRouter()
	// List transactions (index)
	r.Get("/", s.index)
	// New transaction form and create endpoint.
	r.Get("/new", s.new)
	r.Post("/", s.create)
	return r
}

// index handles GET /transactions, renders a paginated list.
func (s *transactionService) index(w http.ResponseWriter, r *http.Request) {
	page := 0
	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p >= 0 {
			page = p
		}
	}
	txs, err := s.repo.Index(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nextURL := fmt.Sprintf("/transactions?page=%d", page+1)
	data := struct {
		Transactions []domain.Transaction
		NextURL      string
		Now          time.Time
	}{
		Transactions: txs,
		NextURL:      nextURL,
		Now:          time.Now(),
	}
	if err := s.tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// new renders the form for creating a new transaction.
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

// create processes the submitted form and creates a new transaction.
func (s *transactionService) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dateStr := r.FormValue("date")
	memo := r.FormValue("memo")
	amtStr := r.FormValue("amount")
	category := r.FormValue("category")
	accountID := r.FormValue("account_id")
	lineItemID := r.FormValue("line_item_id")
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
	tx := domain.Transaction{Date: d, Memo: memo, Amount: amt, Category: category, AccountID: accountID, LineItemID: lineItemID}
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
