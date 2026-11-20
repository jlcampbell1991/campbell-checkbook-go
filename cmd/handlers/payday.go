package handlers

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"

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

// PayDayRoutes returns a chi.Router with pay‑day related endpoints.
func PayDayRoutes(t *template.Template, pdRepo db.PayDayRepository, depRepo db.PayDayDepositRepository, liRepo db.LineItemRepository, acctRepo db.AccountRepository) chi.Router {
    s := &payDayService{repo: pdRepo, depositRepo: depRepo, lineItemRepo: liRepo, accountRepo: acctRepo, tmpl: t}
    r := chi.NewRouter()
    r.Get("/", s.index)
    r.Get("/inactive", s.inactiveIndex)
    r.Get("/new", s.new)
    r.Post("/", s.create)
    r.Get("/{id}", s.show)
    r.Get("/{id}/edit", s.edit)
    r.Post("/{id}/edit", s.update)
    r.Get("/{id}/deactivate", s.deactivate)
    r.Get("/{id}/reactivate", s.reactivate)
    r.Get("/{id}/get-paid", s.getPaid)
    return r
}

func (s *payDayService) index(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    payDays, _ := s.repo.GetAll(ctx)
    data := struct{ PayDays []domain.PayDay }{PayDays: payDays}
    if err := s.tmpl.ExecuteTemplate(w, "pay_days_index.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (s *payDayService) inactiveIndex(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    payDays, _ := s.repo.GetAllInactive(ctx)
    data := struct{ PayDays []domain.PayDay }{PayDays: payDays}
    if err := s.tmpl.ExecuteTemplate(w, "pay_days_inactive.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (s *payDayService) new(w http.ResponseWriter, r *http.Request) {
    // For simplicity, we fetch accounts and line items for selection.
    ctx := r.Context()
    accounts, _ := s.accountRepo.Index(ctx)
    
    data := struct {
        Title    string
        Action   string
        Method   string
        Accounts []domain.Account
    }{
        Title:    "New Pay Day",
        Action:   "/pay-days",
        Method:   "post",
        Accounts: accounts,
    }
    if err := s.tmpl.ExecuteTemplate(w, "pay_days_form.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (s *payDayService) create(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    name := r.FormValue("name")
    totalStr := r.FormValue("total")
    accountId := r.FormValue("accountId")
    if name == "" || totalStr == "" || accountId == "" {
        http.Error(w, "missing required fields", http.StatusBadRequest)
        return
    }
    total, err := strconv.ParseFloat(totalStr, 64)
    if err != nil {
        http.Error(w, "invalid total", http.StatusBadRequest)
        return
    }
    req := &domain.PayDayReq{Name: name, Total: total, AccountId: accountId}
    if _, err := s.repo.Insert(r.Context(), req); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/pay-days", http.StatusSeeOther)
}

func (s *payDayService) show(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")
    pd, err := s.repo.GetById(ctx, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    deposits, _ := s.depositRepo.GetAllByPayDay(ctx, pd.Id)
    data := struct {
        PayDay   domain.PayDay
        Deposits []domain.PayDayDeposit
    }{PayDay: pd, Deposits: deposits}
    if err := s.tmpl.ExecuteTemplate(w, "pay_days_show.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (s *payDayService) edit(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")
    pd, err := s.repo.GetById(ctx, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    accounts, _ := s.accountRepo.Index(ctx)
    data := struct {
        Title    string
        Action   string
        Method   string
        PayDay   domain.PayDay
        Accounts []domain.Account
    }{
        Title:    "Edit Pay Day",
        Action:   fmt.Sprintf("/pay-days/%s/edit", id),
        Method:   "post",
        PayDay:   pd,
        Accounts: accounts,
    }
    if err := s.tmpl.ExecuteTemplate(w, "pay_days_form.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (s *payDayService) update(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    id := chi.URLParam(r, "id")
    name := r.FormValue("name")
    totalStr := r.FormValue("total")
    accountId := r.FormValue("accountId")
    if name == "" || totalStr == "" || accountId == "" {
        http.Error(w, "missing fields", http.StatusBadRequest)
        return
    }
    total, err := strconv.ParseFloat(totalStr, 64)
    if err != nil {
        http.Error(w, "invalid total", http.StatusBadRequest)
        return
    }
    req := &domain.PayDayReq{Name: name, Total: total, AccountId: accountId}
    if _, err := s.repo.Update(r.Context(), req, id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, fmt.Sprintf("/pay-days/%s", id), http.StatusSeeOther)
}

func (s *payDayService) deactivate(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if err := s.repo.Deactivate(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/pay-days", http.StatusSeeOther)
}

func (s *payDayService) reactivate(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if err := s.repo.Reactivate(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/pay-days", http.StatusSeeOther)
}

func (s *payDayService) getPaid(w http.ResponseWriter, r *http.Request) {
    // Stub implementation – in real code this would create transactions.
    http.Redirect(w, r, "/pay-days", http.StatusSeeOther)
}
