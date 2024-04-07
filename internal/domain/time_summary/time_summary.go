package time_summary

type TimeSummary struct {
	MacAddress   string `json:"mac_address"`
	Seconds      int64  `json:"seconds"`
	BreaksJson   []byte `json:"breaks_json"`
	Date         string `json:"date"`
	SecondsStart int64  `json:"seconds_start"`
	SecondsEnd   int64  `json:"seconds_end"`
}
