package main

import (
    "bytes"
    "database/sql"
    _ "github.com/lib/pq"
    "embed"
    "html/template"
    "log"
    "net/http"
    "path/filepath"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

//go:embed views/layouts/*.tmpl views/layouts/*.html views/accounts/*.tmpl
var templateFS embed.FS

var dbConn *sql.DB

func main() {
    // Parse all layout templates
    tmpl, err := template.ParseFS(templateFS, "views/layouts/*.tmpl", "views/layouts/*.html", "views/accounts/*.tmpl")
    if err != nil {
        log.Fatalf("failed to parse templates: %v", err)
    }

    // Initialize DB connection
    dsn := db.InitDB()
    var errOpen error
    dbConn, errOpen = sql.Open("postgres", dsn)
    if errOpen != nil {
        log.Fatalf("failed to open DB: %v", errOpen)
    }
    // Ensure DB is reachable
    if err := dbConn.Ping(); err != nil {
        log.Fatalf("cannot connect to DB: %v", err)
    }

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Serve static assets if present (placeholder paths)
    staticDir := http.Dir(filepath.Join("static"))
    r.Handle("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(staticDir)))
    r.Handle("/javascripts/*", http.StripPrefix("/javascripts/", http.FileServer(staticDir)))

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

    // Account form page handler
    r.Get("/accounts/new", func(w http.ResponseWriter, r *http.Request) {
        var buf bytes.Buffer
        if err := tmpl.ExecuteTemplate(&buf, "add", nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        data := struct {
            Title string
            Body  template.HTML
        }{
            Title: "New Account",
            Body:  template.HTML(buf.String()),
        }
        if err := tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })

    // Handle form submission to create account
    r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "invalid form data", http.StatusBadRequest)
            return
        }
        name := r.FormValue("name")
        if name == "" {
            http.Error(w, "name is required", http.StatusBadRequest)
            return
        }
        if _, err := db.InsertAccount(dbConn, name); err != nil {
            http.Error(w, "failed to create account", http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    log.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
