package handlers

import (
    "html/template"
    "net/http"
    "github.com/go-chi/chi/v5"
)

// RegisterLineItemRoutes registers routes for creating line items.
func RegisterLineItemRoutes(r chi.Router, tmpl *template.Template) {
    r.Get("/accounts/{id}/line-items/new", newLineItemFormHandler(tmpl))
    r.Post("/accounts/{id}/line-items", createLineItemHandler())
}

// newLineItemFormHandler renders the form for creating a new line item.
func newLineItemFormHandler(tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        data := struct {
            AccountID string
        }{AccountID: id}
        if err := tmpl.ExecuteTemplate(w, "lineItemForm", data); err != nil {
            http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
        }
    }
}

// createLineItemHandler processes the form submission.
func createLineItemHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // In a real implementation, parse form and persist.
        // For now, just redirect back to the account page.
        id := chi.URLParam(r, "id")
        http.Redirect(w, r, "/accounts/"+id, http.StatusSeeOther)
    }
}
