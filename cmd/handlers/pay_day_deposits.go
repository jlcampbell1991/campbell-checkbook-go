package handlers

import (
    "context"
    "fmt"
    "html/template"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    db "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/domain"
)

// payDayDepositService holds dependencies for pay day deposit handlers.
type payDayDepositService struct {
    pdRepo       db.PayDayRepository
    depositRepo  db.PayDayDepositRepository
    lineItemRepo db.LineItemRepository
    tmpl         *template.Template
}

// PayDayDepositRoutes creates a chi router for pay day deposit CRUD (new & create).
func PayDayDepositRoutes(t *template.Template, pdRepo db.PayDayRepository, depRepo db.PayDayDepositRepository, liRepo db.LineItemRepository) chi.Router {
    s := &payDayDepositService{pdRepo: pdRepo, depositRepo: depRepo, lineItemRepo: liRepo, tmpl: t}
    r := chi.NewRouter()
    r.Get("/new", s.new)
    r.Post("/", s.create)
    return r
}

// new renders the form for creating a pay day deposit.
func (s *payDayDepositService) new(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // Fetch available pay days.
    payDays, _ := s.pdRepo.GetAll(ctx)
    // Line items – repository currently lacks a List method; use empty slice.
    items := []domain.LineItem{}
    // Optionally we could fetch a specific line item ID from query.
    lineItemID := r.URL.Query().Get("line-item-id")
    data := struct {
        Title         string
        Action        string
        Method        string
        PayDays       []domain.PayDay
        Items         []domain.LineItem
        MaybeItemId   string
        MaybePayDayId string
        MaybeDep      *domain.PayDayDeposit
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

// create processes the form submission to create a pay day deposit.
func (s *payDayDepositService) create(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    amountStr := r.FormValue("amount")
    payDayID := r.FormValue("payDay")
    lineItemID := r.FormValue("lineItem")
    date := r.FormValue("date")
    if amountStr == "" || payDayID == "" || lineItemID == "" {
        http.Error(w, "missing required fields", http.StatusBadRequest)
        return
    }
    amount, err := strconv.ParseFloat(amountStr, 64)
    if err != nil {
        http.Error(w, "invalid amount", http.StatusBadRequest)
        return
    }
    dep := domain.PayDayDeposit{LineItemID: lineItemID, Amount: amount, Date: date}
    // Use Upsert to store the deposit.
    if err := s.depositRepo.Upsert(context.Background(), dep, domain.PayDayId{Value: payDayID}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Redirect to the related line item page.
    redirect := fmt.Sprintf("/line-items/%s", lineItemID)
    http.Redirect(w, r, redirect, http.StatusSeeOther)
}
