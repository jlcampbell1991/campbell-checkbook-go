package main

import (
    "embed"
    "html/template"
    "log"
)

//go:embed views/**/*.tmpl
var tmplFS embed.FS

var tmpl *template.Template

func init() {
    var err error
    tmpl, err = template.New("root").ParseFS(tmplFS, "views/**/*.tmpl")
    if err != nil {
        log.Fatalf("failed to parse templates: %v", err)
    }
}
