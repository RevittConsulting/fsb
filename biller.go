package fsb

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type IBillerDb interface {
	getLastInvoiceDate(ctx context.Context, billingPeriod BillingPeriod) (*time.Time, error)
	createInvoice(ctx context.Context, timeNow time.Time, bill *Bill) error
	getUnpaidInvoices(ctx context.Context) ([]*Invoice, error)
	recordPayment(ctx context.Context, invoice *Invoice, meta []byte, timeNow time.Time) error
}

type Biller struct {
	db IBillerDb
}

func NewBiller(db IBillerDb) *Biller {
	return &Biller{
		db: db,
	}
}

// CreateSchema creates the fsb schema and runs migrations
func (b *Biller) CreateSchema(ctx context.Context, pool *pgxpool.Pool) error {
	return runDBMigrations(ctx, pool)
}

func (b *Biller) CreateInvoice(ctx context.Context, timeNow time.Time, bill *Bill) error {
	lastInvoiceDate, err := b.db.getLastInvoiceDate(ctx, bill.BillingPeriod)
	if err != nil {
		return err
	}

	if lastInvoiceDate == nil {
		return nil
	}

	_, shouldInvoice := CalculateNextInvoiceDate(timeNow, *lastInvoiceDate, bill.BillingPeriod)
	if !shouldInvoice {
		return nil
	}

	return b.db.createInvoice(ctx, timeNow, bill)
}

func (b *Biller) GetUnpaidInvoices(ctx context.Context) ([]*Invoice, error) {
	return b.db.getUnpaidInvoices(ctx)
}

func (b *Biller) RecordPayment(ctx context.Context, invoice *Invoice, meta interface{}, timeNow time.Time) error {
	metaJson, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return b.db.recordPayment(ctx, invoice, metaJson, timeNow)
}
