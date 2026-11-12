package main

import (
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
    "campbell-checkbook/cmd/handlers"
    "campbell-checkbook/internal/db"
)

func main() {
    r := chi.NewRouter()
    // In‑memory repo for demo purposes.
    repo := db.NewMemAccountRepo()
    handlers.RegisterAccounts(r, repo)

    // Simple index handler
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Campbell Checkbook"))
    })

    log.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}
