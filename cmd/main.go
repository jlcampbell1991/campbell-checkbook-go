// Package main implements the HTTP server entry point for the Checkbook application.
package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	dbpkg "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/joho/godotenv"
	"html/template"
)

var tmpl *template.Template

// dbConn holds the database connection used by handlers.
var dbConn *sql.DB

func init() {
	// Parse layout and menu templates. Adjust paths for test execution.
	paths := []string{
		"templates/layouts/application.html",
		"templates/layouts/menu.html",
		"templates/accounts/form.html",
		"templates/accounts/add.html",
		"templates/accounts/index.html",
		"templates/line_items/add.html",
		"templates/accounts/account_items_layout.html",
		"templates/line_items/form.html",
		"templates/accounts/show.html",
	}
	// Resolve paths: if a file does not exist, try with "../" prefix.
	resolved := make([]string, 0, len(paths))
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			resolved = append(resolved, p)
			continue
		}
		alt := "../" + p
		if _, err := os.Stat(alt); err == nil {
			resolved = append(resolved, alt)
			continue
		}
		// If neither exists, keep original to cause a clear error.
		resolved = append(resolved, p)
	}
	tmpl = template.Must(template.ParseFiles(resolved...))
}

type Account struct {
	ID   string
	Name string
}

type LineItemReq struct {
	Name      string
	Balance   float64
	AccountID string
}

type pageData struct {
	Title string
	Body  template.HTML
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Query all accounts from the database
	rows, err := dbConn.Query("SELECT id, name FROM checkbook_accounts ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.ID, &acc.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "accountsIndex", accounts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := pageData{Title: "Checkbook", Body: template.HTML(buf.String())}
	if err := tmpl.ExecuteTemplate(w, "application", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler for displaying the new account form
func newAccountFormHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare data for the form template
	formData := struct {
		Title  string
		Action string
		Method string
		Name   string
	}{
		Title:  "New Account",
		Action: "/accounts",
		Method: "post",
		Name:   "",
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "accountForm", formData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := pageData{Title: "New Account", Body: template.HTML(buf.String())}
	if err := tmpl.ExecuteTemplate(w, "application", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler for processing the submitted new account form
func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	// Generate UUID for the new account
	id := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := dbConn.Exec("INSERT INTO checkbook_accounts (name, id, created_at, updated_at) VALUES ($1, $2, $3, $4)", name, id, now, now)
	if err != nil {
		log.Printf("Failed to insert account: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Redirect to home on success
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handler for displaying the new line item form
func newLineItemFormHandler(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	// Build form data for the template
	formData := struct {
		Action    string
		Method    string
		Name      string
		Balance   float64
		AccountID string
	}{
		Action:    fmt.Sprintf("/accounts/%s/line-items", accountID),
		Method:    "post",
		Name:      "",
		Balance:   0,
		AccountID: accountID,
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "lineItemForm", formData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := pageData{Title: "New Line Item", Body: template.HTML(buf.String())}
	if err := tmpl.ExecuteTemplate(w, "lineItemAdd", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler for creating a new line item
func showAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve account ID from URL
	accountID := chi.URLParam(r, "id")
	// Fetch account details
	var acc Account
	if err := dbConn.QueryRow("SELECT id, name FROM checkbook_accounts WHERE id = $1", accountID).Scan(&acc.ID, &acc.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Fetch active line items for the account
	rows, err := dbConn.Query("SELECT id, name, balance FROM checkbook_line_items WHERE account_id = $1 AND active = true", accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type LineItem struct {
		ID      string
		Name    string
		Balance float64
	}
	var items []LineItem
	for rows.Next() {
		var li LineItem
		if err := rows.Scan(&li.ID, &li.Name, &li.Balance); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, li)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Prepare data for the template
	data := struct {
		Account   Account
		LineItems []LineItem
	}{
		Account:   acc,
		LineItems: items,
	}
	// Render the account show template into a buffer
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "accountShow", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Embed within the application layout
	page := pageData{Title: acc.Name, Body: template.HTML(buf.String())}
	if err := tmpl.ExecuteTemplate(w, "application", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	tx, err := dbConn.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	// Delete line items associated with the account (optional)
	if _, err := tx.Exec("DELETE FROM checkbook_line_items WHERE account_id = $1", accountID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Delete the account itself
	if _, err := tx.Exec("DELETE FROM checkbook_accounts WHERE id = $1", accountID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createLineItemHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	accountID := chi.URLParam(r, "id")
	name := r.FormValue("name")
	balanceStr := r.FormValue("balance")
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		http.Error(w, "Invalid balance", http.StatusBadRequest)
		return
	}
	id := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	_, err = dbConn.Exec(`INSERT INTO checkbook_line_items (name, balance, account_id, id, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, name, balance, accountID, id, true, now, now)
	if err != nil {
		log.Printf("Failed to insert line item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Redirect back to account page or home
	http.Redirect(w, r, fmt.Sprintf("/accounts/%s", accountID), http.StatusSeeOther)
}

func main() {
	router := chi.NewRouter()
	router.Get("/", rootHandler)
	router.Get("/accounts/new", newAccountFormHandler)
	router.Get("/accounts/{id}/line-items/new", newLineItemFormHandler)
	router.Get("/accounts/{id}", showAccountHandler)
	router.Get("/accounts/{id}/delete", deleteAccountHandler)
	router.Post("/accounts/{id}/line-items", createLineItemHandler)
	router.Post("/accounts", createAccountHandler)
	// Serve static files from ./static if present
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Load .env file if present
	_ = godotenv.Load()
	// Initialize DB
	var err error
	dbConn, err = dbpkg.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		portEnv = "8080"
	}
	port := ":" + portEnv
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
