package fsb

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// BillingPeriod represents the period of time for which a bill is generated
// It can be monthly, yearly
type BillingPeriod string

const (
	PeriodMonthly BillingPeriod = "monthly"
	PeriodYearly  BillingPeriod = "yearly"
)

type Bill struct {
	PartyId       uuid.UUID       `db:"party_id" json:"party_id"`
	BillingPeriod BillingPeriod   `db:"billing_period" json:"billing_period"`
	Amount        decimal.Decimal `db:"amount" json:"amount"`
}

func periodToDateTrunc(period BillingPeriod) string {
	switch period {
	case PeriodMonthly:
		return "month"
	case PeriodYearly:
		return "year"
	default:
		return ""
	}
}

type Invoice struct {
	Id        uuid.UUID        `db:"id" json:"id"`
	PartyId   uuid.UUID        `db:"party_id" json:"party_id"`
	Amount    decimal.Decimal  `db:"amount" json:"amount"`
	Meta      *json.RawMessage `db:"meta" json:"meta" sql:"type:jsonb"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
	DeletedAt *time.Time       `db:"deleted_at" json:"deleted_at"`
}

type Payment struct {
	Id        uuid.UUID        `db:"id" json:"id"`
	InvoiceId uuid.UUID        `db:"invoice_id" json:"invoice_id"`
	Meta      *json.RawMessage `db:"meta" json:"meta" sql:"type:jsonb"`
	Success   bool             `db:"success" json:"success"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
	DeletedAt *time.Time       `db:"deleted_at" json:"deleted_at"`
}
