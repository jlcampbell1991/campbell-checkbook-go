package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jlcampbell1991/campbell-checkbook-go/cmd/handlers"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/config"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

func main() {
	// Load DB configuration from environment (and .env file if present)
	cfg := db.Load()
	dsn := db.BuildPostgresURL(cfg)

	// Test DB connection
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open DB connection:", err)
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		log.Fatal("failed to ping DB:", err)
	}
	log.Println("DB connection successful")
	// Initialize pgx pool for repository usage
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("failed to create pgx pool:", err)
	}
	defer pool.Close()

	// Run migrations (no-op if none present)
	db.RunMigrations(dsn)

	// Parse all layout templates
	tmpl, err := template.New("").Funcs(template.FuncMap{"Format": domain.Format, "IsWithdrawal": domain.IsWithdrawal}).ParseGlob("views/**/*.tmpl")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	config.Load()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static assets (CSS and JS)
	// Stylesheets are served from static/stylesheets
	stylesDir := http.Dir(filepath.Join("static", "stylesheets"))
	r.Handle("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(stylesDir)))
	// JavaScripts are served from static/javascripts
	jsDir := http.Dir(filepath.Join("static", "javascripts"))
	r.Handle("/javascripts/*", http.StripPrefix("/javascripts/", http.FileServer(jsDir)))

	// Mount accounts routes (including home page)
	r.Mount("/", handlers.NewAccountsHandler(db.NewAccountRepository(pool), tmpl).Router())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
