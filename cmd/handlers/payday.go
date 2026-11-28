package handlers

import (
    "context"
    "html/template"
    "net/http"
    "github.com/go-chi/chi/v5"
    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

type payDayService struct {
    repo        db.PayDayRepository
    depositRepo db.PayDayDepositRepository
    lineItemRepo db.LineItemRepository
    accountRepo db.AccountRepository
    tmpl        *template.Template
}

func PayDayRoutes(t *template.Template, pdRepo db.PayDayRepository, depRepo db.PayDayDepositRepository, liRepo db.LineItemRepository, acctRepo db.AccountRepository) chi.Router {
    s := &payDayService{repo: pdRepo, depositRepo: depRepo, lineItemRepo: liRepo, accountRepo: acctRepo, tmpl: t}
    r := chi.NewRouter()
    r.Get("/", s.index)
    // Additional routes would be added here.
    return r
}

func (s *payDayService) index(w http.ResponseWriter, r *http.Request) {
    // For now, render a simple template.
    data := struct {
        Title string
    }{Title: "Pay Days"}
    if err := s.tmpl.ExecuteTemplate(w, "pay_days_index.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
