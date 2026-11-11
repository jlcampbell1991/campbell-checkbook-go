package main

import (
	"bytes"
	"html/template"
	"testing"
)

func TestApplicationTemplate(t *testing.T) {
	tmpl := template.Must(template.ParseFiles(
		"views/layouts/menu.html",
		"views/layouts/application.html",
	))
	data := PageData{Title: "Test", Body: template.HTML("<p>Test Body</p>")}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "application", data); err != nil {
		t.Fatalf("template execution failed: %v", err)
	}
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("<h1><a href=\"/\">Checkbook</a></h1>")) {
		t.Errorf("output missing menu header")
	}
	if !bytes.Contains([]byte(out), []byte("<title>Checkbook - Test</title>")) {
		t.Errorf("output missing title")
	}
	if !bytes.Contains([]byte(out), []byte("<p>Test Body</p>")) {
		t.Errorf("output missing body content")
	}
}
