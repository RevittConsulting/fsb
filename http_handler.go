package fsb

import (
	"github.com/RevittConsulting/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(r chi.Router, service *Service) *Handler {
	h := &Handler{
		service: service,
	}
	h.SetupRoutes(r)
	return h
}

func (h *Handler) SetupRoutes(router chi.Router) {
	logger.Log().Info("setting up routes for fsb...")
	router.Group(func(r chi.Router) {
		r.Route("/invoice-payments", func(r chi.Router) {
			r.Get("/", h.GetInvoicePayments)
		})
	})
}

func (h *Handler) GetInvoicePayments(w http.ResponseWriter, r *http.Request) {
	invoicePayments, err := h.service.GetInvoicePayments(r.Context())
	if err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, invoicePayments)
}
