package handlers

import (
    "html/template"
    "log"
    "net/http"
)

// getNewAccountForm renders a simple HTML form for creating a new account.
func GetNewAccountForm(tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Simple form HTML; could be a separate template in the future.
        formHTML := `<form action="/accounts" method="post">
    Name: <input type="text" name="name"><br>
    <input type="submit" value="Submit">
</form>`
        data := struct {
            Title string
            Body  template.HTML
        }{
            Title: "New Account",
            Body:  template.HTML(formHTML),
        }
        if err := tmpl.ExecuteTemplate(w, "application.tmpl", data); err != nil {
            log.Printf("template execution error: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}

// postAccount processes the submitted account form.
func PostAccount(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "invalid form", http.StatusBadRequest)
        return
    }
    name := r.FormValue("name")
    // In a full implementation, we would persist the account.
    log.Printf("Created account with name: %s", name)
    // Redirect to home page after creation.
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
