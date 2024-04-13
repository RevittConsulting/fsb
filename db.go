package fsb

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Db struct {
	pool *pgxpool.Pool
}

func NewDb(pool *pgxpool.Pool) *Db {
	return &Db{
		pool: pool,
	}
}

func (db *Db) getLastInvoiceDate(ctx context.Context, billingPeriod BillingPeriod) (*time.Time, error) {
	sql := `select party_id, max(created_at) as last_invoice_date
			from fsa.invoices
			where created_at >= date_trunc($1, current_date)
			group by party_id
			having max(created_at) is null or max(created_at) is not null;`

	var lastInvoiceDate time.Time
	err := db.pool.QueryRow(ctx, sql, periodToDateTrunc(billingPeriod)).Scan(&lastInvoiceDate)
	if err != nil {
		return nil, err
	}

	return &lastInvoiceDate, nil
}

func (db *Db) createInvoice(ctx context.Context, timeNow time.Time, bill *Bill) error {
	sql := `insert into fsa.invoices (party_id, amount, created_at)
			values ($1, $2, $3);`

	tx, err := txBegin(ctx, db.pool)
	if err != nil {
		return err
	}
	defer txDefer(tx, ctx)

	_, err = tx.Exec(ctx, sql, bill.PartyId, bill.Amount, timeNow)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) getUnpaidInvoices(ctx context.Context) ([]*Invoice, error) {
	sql := `select i.id, i.party_id, i.meta, i.amount
			from chd.invoices i
			left join chd.invoice_payments ip on i.id = ip.invoice_id
			where ip.id is null;`

	rows, err := db.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	var invoices []*Invoice
	for rows.Next() {
		invoice := &Invoice{}
		err = rows.Scan(&invoice.Id, &invoice.PartyId, &invoice.Meta, &invoice.Amount)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (db *Db) recordPayment(ctx context.Context, invoice *Invoice, meta []byte, timeNow time.Time) error {
	sql := `insert into chd.invoice_payments (meta, invoice_id, created_at)
			values ($1, $2, $3);`

	tx, err := txBegin(ctx, db.pool)
	if err != nil {
		return err
	}
	defer txDefer(tx, ctx)

	if _, err = tx.Exec(ctx, sql, meta, invoice.Id, timeNow); err != nil {
		return err
	}

	err = txCommit(tx, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) getInvoicePayments(ctx context.Context) ([]*Payment, error) {
	sql := `select ip.id, ip.invoice_id, ip.meta, ip.created_at
			from chd.invoice_payments ip;`

	rows, err := db.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	var payments []*Payment
	for rows.Next() {
		payment := &Payment{}
		err = rows.Scan(&payment.Id, &payment.InvoiceId, &payment.Meta, &payment.CreatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}
