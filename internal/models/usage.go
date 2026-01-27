package models

// MonthSummary represents the monthly summary data
type MonthSummary struct {
	ID                                  string  `json:"_id"`
	ActualCostsUntilLastMeterReadingDate float64 `json:"actualCostsUntilLastMeterReadingDate"`
	ExpectedCostsUntilLastMeterReadingDate float64 `json:"expectedCostsUntilLastMeterReadingDate"`
	ExpectedCosts                        float64 `json:"expectedCosts"`
	LastMeterReadingDate                 string  `json:"lastMeterReadingDate"`
	MeterReadingDayCompleteness          float64 `json:"meterReadingDayCompleteness"`
	GasExcluded                          bool    `json:"gasExcluded"`
}

// MonthSummaryResponse represents the API response
type MonthSummaryResponse struct {
	MonthSummary *MonthSummary `json:"monthSummary"`
}

// UsageItem represents a single usage entry
type UsageItem struct {
	Date  string  `json:"date"`
	From  string  `json:"from"`
	Till  string  `json:"till"`
	Usage float64 `json:"usage"`
	Costs float64 `json:"costs"`
	Unit  string  `json:"unit"`
}

// EnergyCategory represents usage for a category (electricity, gas, feedIn)
type EnergyCategory struct {
	UsageTotal float64     `json:"usageTotal"`
	CostsTotal float64     `json:"costsTotal"`
	Unit       string      `json:"unit"`
	Items      []UsageItem `json:"items"`
}

// PeriodUsageAndCosts represents usage and costs for a period
type PeriodUsageAndCosts struct {
	ID          string          `json:"_id"`
	Gas         *EnergyCategory `json:"gas"`
	Electricity *EnergyCategory `json:"electricity"`
	FeedIn      *EnergyCategory `json:"feedIn"`
}

// PeriodUsageAndCostsResponse represents the API response
type PeriodUsageAndCostsResponse struct {
	PeriodUsageAndCosts *PeriodUsageAndCosts `json:"periodUsageAndCosts"`
}

// Invoice represents a single invoice
type Invoice struct {
	ID                string  `json:"id"`
	InvoiceDate       string  `json:"invoiceDate"`
	StartDate         string  `json:"startDate"`
	PeriodDescription string  `json:"periodDescription"`
	TotalAmount       float64 `json:"totalAmount"`
}

// Invoices represents the invoices response
type Invoices struct {
	AllInvoices           []Invoice `json:"allInvoices"`
	PreviousPeriodInvoice *Invoice  `json:"previousPeriodInvoice"`
	CurrentPeriodInvoice  *Invoice  `json:"currentPeriodInvoice"`
	UpcomingPeriodInvoice *Invoice  `json:"upcomingPeriodInvoice"`
}

// InvoicesResponse represents the API response
type InvoicesResponse struct {
	Invoices *Invoices `json:"invoices"`
}
