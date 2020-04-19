package etracking

import "time"

type TaxResult struct {
	Available bool

	Barcode           string
	CustomID          string
	Recipient         string
	ReceivingLocation string

	ImportTax     int64 // unit is satang
	ExciseTax     int64 // unit is satang
	InteriorTax   int64 // unit is satang
	ValueAddedTax int64 // unit is satang
	OtherFee      int64 // unit is satang
	TotalTax      int64 // unit is satang

	Steps []CustomStepEntry
}

type CustomStepEntry struct {
	Step string
	Time time.Time
}
