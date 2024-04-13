package fsb

import (
	"context"
	"github.com/RevittConsulting/logger"
	"go.uber.org/zap"
	"time"
)

type ITimeProvider interface {
	Now() time.Time
}

type IPaymentProvider interface {
	CreatePayment(ctx context.Context) (interface{}, error)
}

type ISubscriptionProvider interface {
	GetBills(ctx context.Context) ([]*Bill, error)
}

type Scheduler struct {
	biller  *Biller
	tp      ITimeProvider
	payment IPaymentProvider
	subs    ISubscriptionProvider
}

func NewScheduler(tp ITimeProvider, biller *Biller, payment IPaymentProvider, subs ISubscriptionProvider) *Scheduler {
	return &Scheduler{
		tp:      tp,
		biller:  biller,
		payment: payment,
		subs:    subs,
	}
}

func (s *Scheduler) RunInvoiceJob(ctx context.Context) {
	bills, err := s.subs.GetBills(ctx)
	if err != nil {
		logger.Log().Error("error getting bills: %v", zap.Error(err))
	}

	for _, bill := range bills {
		err := s.biller.CreateInvoice(ctx, s.tp.Now(), bill)
		if err != nil {
			logger.Log().Error("error creating invoice: %v", zap.Error(err))
		}
	}
}

func (s *Scheduler) RunPaymentJob(ctx context.Context) {
	invoices, err := s.biller.GetUnpaidInvoices(ctx)
	if err != nil {
		logger.Log().Error("error getting unpaid invoices: %v", zap.Error(err))
	}

	for _, invoice := range invoices {
		meta, err := s.payment.CreatePayment(ctx)
		if err != nil {
			logger.Log().Error("error creating payment: %v", zap.Error(err))
			continue
		}

		if err = s.biller.RecordPayment(ctx, invoice, meta, s.tp.Now()); err != nil {
			logger.Log().Error("error recording payment: %v", zap.Error(err))
		}
	}
}
