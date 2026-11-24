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

type payDayDepositService struct {
    pdRepo     db.PayDayRepository
    depositRepo db.PayDayDepositRepository
    liRepo     db.LineItemRepository
    tmpl       *template.Template
}

// PayDayDepositRoutes returns a chi.Router handling Pay Day Deposit creation.
func PayDayDepositRoutes(t *template.Template, pdRepo db.PayDayRepository, depositRepo db.PayDayDepositRepository, liRepo db.LineItemRepository) chi.Router {
    s := &payDayDepositService{pdRepo: pdRepo, depositRepo: depositRepo, liRepo: liRepo, tmpl: t}
    r := chi.NewRouter()
    r.Get("/new", s.new)
    r.Post("/", s.create)
    return r
}

func (s *payDayDepositService) new(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // Optional line-item-id query param
    lineItemID := r.URL.Query().Get("line-item-id")
    // Fetch PayDays for dropdown
    payDays, _ := s.pdRepo.GetAll(ctx)
    // For simplicity, we pass empty items slice; can be extended later.
    items := []domain.LineItem{}
    data := struct {
        Title        string
        Action       string
        Method       string
        PayDays      []domain.PayDay
        Items        []domain.LineItem
        MaybeItemId  string
    }{
        Title:       "New Pay Day Deposit",
        Action:      "/pay-day-deposits",
        Method:      "post",
        PayDays:     payDays,
        Items:       items,
        MaybeItemId: lineItemID,
    }
    if err := s.tmpl.ExecuteTemplate(w, "pay_day_deposits_add.tmpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (s *payDayDepositService) create(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    amountStr := r.FormValue("amount")
    payDayID := r.FormValue("payDay")
    lineItemID := r.FormValue("lineItem")
    date := r.FormValue("date")
    amt, err := strconv.ParseFloat(amountStr, 64)
    if err != nil {
        http.Error(w, "invalid amount", http.StatusBadRequest)
        return
    }
    dep := domain.PayDayDeposit{LineItemID: lineItemID, Amount: amt, Date: date}
    // Insert deposit via repository; ignore returned ID for now.
    if err := s.depositRepo.Upsert(r.Context(), dep, domain.PayDayId{Value: payDayID}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Redirect to line‑item page.
    http.Redirect(w, r, fmt.Sprintf("/line-items/%s", lineItemID), http.StatusSeeOther)
}
