package time_summary

import "time"

type TimeSummary struct {
	MacAddress   string
	Seconds      int64
	BreaksJson   []byte
	Date         time.Time
	SecondsStart int64
	SecondsEnd   int64
}
