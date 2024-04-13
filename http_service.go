package fsb

import "context"

type IHttpDb interface {
	getInvoicePayments(ctx context.Context) ([]*Payment, error)
}

type Service struct {
	db IHttpDb
}

func NewService(db IHttpDb) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetInvoicePayments(ctx context.Context) ([]*Payment, error) {
	return s.db.getInvoicePayments(ctx)
}
