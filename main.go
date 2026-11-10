package main

import (
    "embed"
    "html/template"
    "log"
    "net/http"
    "path/filepath"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

//go:embed views/layouts/*.tmpl views/layouts/*.html
var templateFS embed.FS

func main() {
    // Parse all layout templates
    tmpl, err := template.ParseFS(templateFS, "views/layouts/*.tmpl", "views/layouts/*.html")
    if err != nil {
        log.Fatalf("failed to parse templates: %v", err)
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

    log.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
