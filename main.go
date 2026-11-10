package main

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl *template.Template

func init() {
	// Parse layout and menu templates
	tmpl = template.Must(template.ParseFiles(
		"templates/layouts/application.html",
		"templates/layouts/menu.html",
		"templates/accounts/form.html",
		"templates/accounts/add.html",
	))
}

type pageData struct {
	Title string
	Body  template.HTML
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		Title: "Checkbook",
		Body:  "Welcome to Checkbook",
	}
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
	// Simple placeholder Account struct (could be expanded later)
	type Account struct {
		Name string
	}
	account := Account{Name: name}
	msg := "Account created: " + account.Name
	page := pageData{Title: "Account Created", Body: template.HTML(msg)}
	if err := tmpl.ExecuteTemplate(w, "application", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	router := chi.NewRouter()
	router.Get("/", rootHandler)
	router.Get("/accounts/new", newAccountFormHandler)
	router.Post("/accounts", createAccountHandler)
	// Serve static files from ./static if present
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Load .env file if present
	_ = godotenv.Load()
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
