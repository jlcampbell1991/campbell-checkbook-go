package main

import (
    _ "embed"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
    "github.com/jlcampbell1991/campbell-checkbook-go/cmd/handlers"
    "github.com/joho/godotenv"
)


var accounts = map[int]domain.Account{
    1: {ID: 1, Name: "Checking", Balance: 1000.00},
    2: {ID: 2, Name: "Savings", Balance: 2500.50},
}

func main() {
    // Templates are parsed in init()

    // Load environment variables and run DB migrations
    _ = godotenv.Load()
    cfg := db.PostgresConfig{
        User:     os.Getenv("POSTGRES_USER"),
        Password: os.Getenv("POSTGRES_PASSWORD"),
        Host:     os.Getenv("POSTGRES_HOST"),
        DB:       os.Getenv("POSTGRES_DB"),
        Port:     os.Getenv("POSTGRES_PORT"),
    }
    dsn := db.BuildPostgresURL(cfg)
    db.RunMigrations(dsn)

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Serve static assets if present (placeholder paths)
    staticDir := http.Dir(filepath.Join("static"))
    r.Handle("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(staticDir)))
    r.Handle("/javascripts/*", http.StripPrefix("/javascripts/", http.FileServer(staticDir)))

	// Account show handler
	r.Get("/accounts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid account id", http.StatusBadRequest)
			return
		}
		account, ok := accounts[id]
		if !ok {
			http.NotFound(w, r)
			return
		}
		body := template.HTML(
			"<h1>Account " + account.Name + "</h1>" +
				"<p>ID: " + strconv.Itoa(account.ID) + "</p>" +
				fmt.Sprintf("<p>Balance: $%.2f</p>", account.Balance) +
				"<form action=\"/accounts/"+strconv.Itoa(account.ID)+"/delete\" method=\"POST\"><button type=\"submit\">Delete</button></form>")
		data := struct {
			Title string
			Body  template.HTML
		}{
			Title: "Account " + account.Name,
			Body:  body,
		}
		if err := tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Account delete handler
	r.Post("/accounts/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid account id", http.StatusBadRequest)
			return
		}
		if _, ok := accounts[id]; !ok {
			http.NotFound(w, r)
			return
		}
		delete(accounts, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	})
    // Transaction routes
    handlers.RegisterTransactionRoutes(r, tmpl)


    // Home page handler
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        data := struct {
            Title string
            Body  template.HTML
        }{
            Title: "Home",
            Body:  template.HTML("<p>Welcome to Campbell Checkbook.</p>"),
        }
        // Execute the application layout template
        if err := tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server listening on :%s", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
