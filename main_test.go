package main

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func setupRouter() (*chi.Mux, *template.Template) {
	r := chi.NewRouter()
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "application.html"),
		filepath.Join("views", "layouts", "menu.html"),
		filepath.Join("views", "accounts", "add.html"),
	))
	// Home (not needed for tests but keep consistent)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct{ Title string }{Title: "Home"}
		_ = tmpl.ExecuteTemplate(w, "application", data)
	})
	// Account routes
	r.Get("/accounts/new", func(w http.ResponseWriter, r *http.Request) {
		data := struct{ Title string }{Title: "New Account"}
		_ = tmpl.ExecuteTemplate(w, "add", data)
	})
	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		_ = r.FormValue("name")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	return r, tmpl
}

func TestGetNewAccountForm(t *testing.T) {
	r, _ := setupRouter()
	req := httptest.NewRequest("GET", "/accounts/new", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	body := w.Body.String()
	if body == "" || !contains(body, "New Account") {
		t.Fatalf("response does not contain expected content")
	}
}

func TestPostCreateAccountRedirect(t *testing.T) {
	r, _ := setupRouter()
	form := "name=TestAccount"
	req := httptest.NewRequest("POST", "/accounts", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected redirect status, got %d", resp.StatusCode)
	}
	loc, err := resp.Location()
	if err != nil {
		t.Fatalf("expected location header: %v", err)
	}
	if loc.Path != "/" {
		t.Fatalf("expected redirect to '/', got %s", loc.Path)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
