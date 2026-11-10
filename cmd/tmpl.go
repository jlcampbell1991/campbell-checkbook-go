package main

import (
    "html/template"
    "log"
)

var tmpl *template.Template

func init() {
    var err error
    tmpl = template.New("root")
    tmpl, err = tmpl.ParseGlob("../views/*.tmpl")
    if err != nil {
        log.Fatalf("failed to parse templates: %v", err)
    }
    tmpl, err = tmpl.ParseGlob("../views/layouts/*.tmpl")
    if err != nil {
        log.Fatalf("failed to parse templates: %v", err)
    }
    if err != nil {
        log.Fatalf("failed to parse templates: %v", err)
    }
}
