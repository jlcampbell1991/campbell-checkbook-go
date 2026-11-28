package handlers

import (
        "fmt"
    "html/template"
    "net/http"
    "strconv"
    "path/filepath"

    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)


func IsWithdrawal(amount float64) bool { return amount < 0 }
func CurrencyFromDouble(v float64) string { return fmt.Sprintf("%.2f", v) }

// TransactionsIndex renders the transactions list page.
func TransactionsIndex(w http.ResponseWriter, r *http.Request) {
    // Parse page query parameter, default to 0.
    page := 0
    if p := r.URL.Query().Get("page"); p != "" {
        if v, err := strconv.Atoi(p); err == nil && v >= 0 {
            page = v
        }
    }

    // Retrieve transactions (stub implementation).
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

    // Parse templates including layout and transaction templates.
    // Parse required templates from the file system.
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
