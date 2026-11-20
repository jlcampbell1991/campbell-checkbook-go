package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"github.com/jlcampbell1991/campbell-checkbook-go/utils"
)

// tmpl holds the parsed templates for backward compatibility.
var tmpl *template.Template

// SetTemplate makes the parsed templates available to legacy code.
func SetTemplate(t *template.Template) {
	tmpl = t
}

// transactionService holds dependencies for transaction handlers.
type transactionService struct {
	repo db.TransactionRepository
	tmpl *template.Template
}

// TransactionRoutes creates a chi router for transaction CRUD.
func TransactionRoutes(t *template.Template, repo db.TransactionRepository) chi.Router {
	svc := &transactionService{repo: repo, tmpl: t}
	r := chi.NewRouter()
	// List transactions (index)
	r.Get("/", svc.index)
	// Show form for new transaction
	r.Get("/new", svc.new)
	// Handle creation
	r.Post("/", svc.create)
	return r
}

// index handles GET /transactions (list).
func (s *transactionService) index(w http.ResponseWriter, r *http.Request) {
	page := 0
	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p >= 0 {
			page = p
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

	s.tmpl = s.tmpl.Funcs(template.FuncMap{"isWithdrawal": utils.IsWithdrawal})
	if err := s.tmpl.ExecuteTemplate(w, "transactions_index.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// new renders the form for creating a transaction.
func (s *transactionService) new(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New Transaction Form"))
}

// create handles POST form submission to create a transaction.
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

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amtStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	tx := domain.Transaction{Date: date, Memo: memo, Amount: amount, Category: category, AccountID: accountID, LineItemID: lineItemID}
	if _, err := s.repo.Create(context.Background(), tx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirect := "/"
	if lineItemID != "" {
		redirect = fmt.Sprintf("/line-items/%s", lineItemID)
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
