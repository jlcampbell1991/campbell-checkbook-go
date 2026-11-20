package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"github.com/jlcampbell1991/campbell-checkbook-go/utils"
)

// tmpl holds the parsed templates and is set from main.
var tmpl *template.Template

// SetTemplate makes the parsed templates available to this handler.
func SetTemplate(t *template.Template) {
	tmpl = t
}

// TransactionsIndex handles the GET /transactions page.
func TransactionsIndex(w http.ResponseWriter, r *http.Request) {
	// Parse page query parameter, default 0.
	page := 0
	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p >= 0 {
			page = p
		}
	}

	// Use stub repository for now.
	repo := db.NewStubTransactionRepo()
	transactions, err := repo.Index(page)
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

	// Ensure template functions are available.
	tmpl = tmpl.Funcs(template.FuncMap{"isWithdrawal": utils.IsWithdrawal})

	if err := tmpl.ExecuteTemplate(w, "transactions_index.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
