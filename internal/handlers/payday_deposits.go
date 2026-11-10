package handlers

import (
    "html/template"
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
)

type PayDayDepositFormData struct {
    Title  string
    Action string
    Method string
    // For simplicity, we won't populate dynamic data
}

// RegisterRoutes registers pay day deposit routes onto the router.
func RegisterRoutes(r chi.Router, tmpl *template.Template) {
    r.Get("/pay-day-deposits/add", func(w http.ResponseWriter, r *http.Request) {
        data := PayDayDepositFormData{Title: "New Pay Day Deposit", Action: "/pay-day-deposits", Method: "post"}
        // Render the form template within the application layout
        // Assuming a layout template named "application.tmpl" that expects .Title and .Body
        // We'll embed the form directly as body HTML.
        // For simplicity, we execute the form template and capture output.
        var bodyHTML string
        // Execute form template to string
        // Using a temporary ResponseRecorder
        rr := &responseRecorder{header: http.Header{}}
        if err := tmpl.ExecuteTemplate(rr, "pay_day_deposit_form", data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        bodyHTML = rr.body
        pageData := struct {
            Title string
            Body  template.HTML
        }{Title: data.Title, Body: template.HTML(bodyHTML)}
        if err := tmpl.ExecuteTemplate(w, "application.tmpl", pageData); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })

    r.Post("/pay-day-deposits", func(w http.ResponseWriter, r *http.Request) {
        // Simple stub: parse form values and log them
        if err := r.ParseForm(); err != nil {
            http.Error(w, "invalid form", http.StatusBadRequest)
            return
        }
        log.Printf("Received PayDayDeposit form: %v", r.PostForm)
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })
}

// responseRecorder is a minimal http.ResponseWriter to capture template output.
type responseRecorder struct {
    header http.Header
    body   string
    status int
}

func (rr *responseRecorder) Header() http.Header { return rr.header }
func (rr *responseRecorder) Write(b []byte) (int, error) {
    rr.body += string(b)
    return len(b), nil
}
func (rr *responseRecorder) WriteHeader(statusCode int) { rr.status = statusCode }
