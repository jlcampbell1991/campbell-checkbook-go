package handlers

import (
        "fmt"
    "html/template"
    "net/http"
    "path/filepath"
    "strconv"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

func IsWithdrawal(amount float64) bool { return amount < 0 }
func CurrencyFromDouble(v float64) string { return fmt.Sprintf("%.2f", v) }

// TransactionsIndex renders the transactions list page (legacy handler).
func TransactionsIndex(w http.ResponseWriter, r *http.Request) {
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

// transactionService provides handlers for transaction CRUD.
type transactionService struct {
    repo db.TransactionRepository
    tmpl *template.Template
}

// TransactionRoutes returns a chi.Router handling transaction endpoints.
func TransactionRoutes(t *template.Template, repo db.TransactionRepository) chi.Router {
    s := &transactionService{repo: repo, tmpl: t}
    r := chi.NewRouter()
    r.Get("/", s.index)
    r.Get("/new", s.new)
    r.Post("/", s.create)
    // Additional routes (edit, delete, etc.) can be added here.
    return r
}

// index renders the list of transactions using the injected repository.
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

// create processes the form submission to add a transaction.
func (s *transactionService) create(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    dateStr := r.FormValue("date")
    memo := r.FormValue("memo")
    amountStr := r.FormValue("amount")
    category := r.FormValue("category")
    accountID := r.FormValue("account_id")
    lineItemID := r.FormValue("line_item_id")
    if dateStr == "" || memo == "" || amountStr == "" {
        http.Error(w, "missing required fields", http.StatusBadRequest)
        return
    }
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        http.Error(w, "invalid date", http.StatusBadRequest)
        return
    }
    amount, err := strconv.ParseFloat(amountStr, 64)
    if err != nil {
        http.Error(w, "invalid amount", http.StatusBadRequest)
        return
    }
    tx := domain.Transaction{Date: date, Memo: memo, Amount: amount, Category: category, AccountID: accountID, LineItemID: lineItemID}
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
