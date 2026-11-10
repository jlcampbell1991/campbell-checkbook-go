package handlers

import (
    "html/template"
    "net/http"
    "github.com/go-chi/chi/v5"
    "fmt"
    "strings"
)

// RegisterLineItemRoutes registers routes for creating and showing line items.
func RegisterLineItemRoutes(r chi.Router, tmpl *template.Template) {
    // Show a line item
    r.Get("/line-items/{id}", showLineItemHandler(tmpl))
    // New line item form
    r.Get("/accounts/{id}/line-items/new", newLineItemFormHandler(tmpl))
    // Create line item
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

// showLineItemHandler renders a line item show page.
func showLineItemHandler(tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        // Dummy data for line item
        data := struct {
            LineItem struct {
                ID    string
                Name  string
                Balance float64
            }
        }{}
        data.LineItem.ID = id
        data.LineItem.Name = fmt.Sprintf("Line Item %s", id)
        data.LineItem.Balance = 123.45

        // Render lineItemShow into buffer then layout
        var buf strings.Builder
        if err := tmpl.ExecuteTemplate(&buf, "lineItemShow", data); err != nil {
            http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        pageData := struct {
            Title string
            Body  template.HTML
        }{
            Title: fmt.Sprintf("Line Item %s", id),
            Body:  template.HTML(buf.String()),
        }
        if err := tmpl.ExecuteTemplate(w, "application.tmpl", pageData); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}
