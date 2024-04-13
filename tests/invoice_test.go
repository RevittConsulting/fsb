package tests

import (
	"testing"
	"time"

	"github.com/RevittConsulting/fsb"
)

func Test_CalculateNextInvoiceDate(t *testing.T) {
	scenarios := map[string]struct {
		timeNow           time.Time
		lastInvoiceDate   time.Time
		billingPeriod     fsb.BillingPeriod
		expectedDate      time.Time
		expectedToInvoice bool
	}{
		"EndOfMonthNonLeapYear ShouldInvoice": {
			timeNow:           time.Date(2023, 3, 01, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodMonthly,
			expectedDate:      time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: true,
		},
		"EndOfMonthNonLeapYear ShouldNotInvoice": {
			timeNow:           time.Date(2023, 2, 26, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodMonthly,
			expectedDate:      time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: false,
		},
		"EndOfMonthLeapYear ShouldInvoice": {
			timeNow:           time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodMonthly,
			expectedDate:      time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: true,
		},
		"EndOfMonthLeapYear ShouldNotInvoice": {
			timeNow:           time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodMonthly,
			expectedDate:      time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: false,
		},
		"RegularDay ShouldInvoice": {
			timeNow:           time.Date(2023, 5, 16, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodMonthly,
			expectedDate:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: true,
		},
		"RegularDay ShouldNotInvoice": {
			timeNow:           time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodMonthly,
			expectedDate:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: false,
		},
		"AnnualBillingLeapYear ShouldInvoice": {
			timeNow:           time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodYearly,
			expectedDate:      time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: true,
		},
		"AnnualBillingLeapYear ShouldNotInvoice": {
			timeNow:           time.Date(2021, 1, 28, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodYearly,
			expectedDate:      time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: false,
		},
		"AnnualBillingNonLeapYear ShouldInvoice": {
			timeNow:           time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodYearly,
			expectedDate:      time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: true,
		},
		"AnnualBillingNonLeapYear ShouldNotInvoice": {
			timeNow:           time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			lastInvoiceDate:   time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			billingPeriod:     fsb.PeriodYearly,
			expectedDate:      time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
			expectedToInvoice: false,
		},
	}

	for name, scenario := range scenarios {
		t.Run(name, func(t *testing.T) {
			actualDate, shouldInvoice := fsb.CalculateNextInvoiceDate(scenario.lastInvoiceDate, scenario.timeNow, scenario.billingPeriod)
			if actualDate == nil || *actualDate != scenario.expectedDate {
				t.Errorf("Scenario %s failed: expected %v, got %v", name, scenario.expectedDate, actualDate)
			}
			if shouldInvoice != scenario.expectedToInvoice {
				t.Errorf("Scenario %s failed: expected to invoice %v, got %v", name, scenario.expectedToInvoice, shouldInvoice)
			}
		})
	}
}
