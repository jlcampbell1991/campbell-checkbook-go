package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var tmplFS embed.FS

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFS(tmplFS, "templates/*.html")
	if err != nil {
		panic(err)
	}
}

type LineItem struct {
	ID      string
	Name    string
	Balance float64
}

type LineItemShowData struct {
	Item      LineItem
	DeleteURL string
}

func lineItemShowHandler(w http.ResponseWriter, r *http.Request) {
	// For demo purposes, we use a static line item.
	item := LineItem{ID: "li-123", Name: "Test LineItem", Balance: 42.0}
	data := LineItemShowData{Item: item, DeleteURL: fmt.Sprintf("/line-items/%s/delete", item.ID)}
	err := tmpl.ExecuteTemplate(w, "lineItemShow", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/line-items/", lineItemShowHandler)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
