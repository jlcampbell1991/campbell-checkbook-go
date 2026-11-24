package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	"github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
	"html/template"
	"net/http"
	"strconv"
)

type payDayService struct {
	repo         db.PayDayRepository
	depositRepo  db.PayDayDepositRepository
	lineItemRepo db.LineItemRepository
	accountRepo  db.AccountRepository
	tmpl         *template.Template
}

// PayDayRoutes returns a chi.Router with pay‑day endpoints.
func PayDayRoutes(t *template.Template, pdRepo db.PayDayRepository, depRepo db.PayDayDepositRepository, liRepo db.LineItemRepository, acctRepo db.AccountRepository) chi.Router {
	s := &payDayService{repo: pdRepo, depositRepo: depRepo, lineItemRepo: liRepo, accountRepo: acctRepo, tmpl: t}
	r := chi.NewRouter()
	r.Get("/", s.index)
	r.Get("/inactive-pay-days", s.inactiveIndex)
	r.Get("/new", s.new)
	r.Post("/", s.create)
	r.Get("/{id}", s.show)
	r.Get("/{id}/edit", s.edit)
	r.Post("/{id}/edit", s.update)
	r.Get("/{id}/deactivate", s.deactivate)
	r.Get("/{id}/reactivate", s.reactivate)
	// Additional routes omitted for brevity.
	return r
}

func (s *payDayService) index(w http.ResponseWriter, r *http.Request) {
	// Stub implementation: render empty list.
	data := struct{ PayDays []domain.PayDay }{PayDays: []domain.PayDay{}}
	if err := s.tmpl.ExecuteTemplate(w, "pay_days_index.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *payDayService) inactiveIndex(w http.ResponseWriter, r *http.Request) {
	data := struct{ PayDays []domain.PayDay }{PayDays: []domain.PayDay{}}
	if err := s.tmpl.ExecuteTemplate(w, "pay_days_inactive_index.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *payDayService) new(w http.ResponseWriter, r *http.Request) {
	// Render form for creating a new PayDay.
	// For simplicity, provide empty struct.
	data := struct {
		Title   string
		Action  string
		Method  string
		MaybePD *domain.PayDayReq
	}{
		Title:  "New Pay Day",
		Action: "/pay-days",
		Method: "post",
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
	accountId := r.FormValue("account")
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
	id := chi.URLParam(r, "id")
	pd, err := s.repo.GetById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// For simplicity, no deposits loaded.
	data := struct {
		PayDay   domain.PayDay
		Deposits []domain.PayDayDeposit
	}{PayDay: pd, Deposits: []domain.PayDayDeposit{}}
	if err := s.tmpl.ExecuteTemplate(w, "pay_days_show.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *payDayService) edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pd, err := s.repo.GetById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Convert to PayDayReq for form population.
	req := &domain.PayDayReq{Name: pd.Name, Total: pd.Total, AccountId: pd.AccountId}
	data := struct {
		Title   string
		Action  string
		Method  string
		MaybePD *domain.PayDayReq
	}{
		Title:   "Edit Pay Day",
		Action:  fmt.Sprintf("/pay-days/%s/edit", id),
		Method:  "post",
		MaybePD: req,
	}
	if err := s.tmpl.ExecuteTemplate(w, "pay_days_form.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *payDayService) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	totalStr := r.FormValue("total")
	accountId := r.FormValue("account")
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
	http.Redirect(w, r, "/pay-days/inactive-pay-days", http.StatusSeeOther)
}
