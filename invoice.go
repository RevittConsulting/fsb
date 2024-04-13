package fsb

import (
	"time"
)

func CalculateNextInvoiceDate(lastInvoiceDate, now time.Time, period BillingPeriod) (*time.Time, bool) {
	var nextBillingDate time.Time

	switch period {
	case PeriodMonthly:
		_, _, day := lastInvoiceDate.Date()
		nextMonth := lastInvoiceDate.AddDate(0, 1, 0)
		daysInNextMonth := daysInMonth(lastInvoiceDate.Year(), lastInvoiceDate.Month())

		if nextMonth.Day() < day {
			dif := day - daysInNextMonth
			nextBillingDate = lastInvoiceDate.AddDate(0, 1, -dif)
		} else {
			nextBillingDate = nextMonth
		}

	case PeriodYearly:
		nextYear := lastInvoiceDate.AddDate(1, 0, 0)

		if lastInvoiceDate.Month() == time.February && lastInvoiceDate.Day() == 29 && !isLeapYear(nextYear.Year()) {
			nextBillingDate = time.Date(nextYear.Year(), time.February, 28, 0, 0, 0, 0, time.UTC)
		} else {
			nextBillingDate = nextYear
		}
	}

	nowTrunc := now.Truncate(24 * time.Hour)
	nextBillingDateTrunc := nextBillingDate.Truncate(24 * time.Hour)

	if nextBillingDateTrunc.Before(nowTrunc) {
		return &nextBillingDate, true
	}

	return &nextBillingDate, false
}

func daysInMonth(year int, month time.Month) int {
	start := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	days := end.Add(-24 * time.Hour)
	return days.Day()
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
