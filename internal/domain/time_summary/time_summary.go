package time_summary

// Break ???
type Break struct {
	SecondsStart int64 `json:"seconds_start"`
	SecondsEnd   int64 `json:"seconds_end"`
}

// TimeSummary ???
type TimeSummary struct {
	MacAddress   string   `json:"mac_address"`
	Seconds      int64    `json:"seconds"`
	Breaks       []*Break `json:"breaks"`
	Date         string   `json:"date"`
	SecondsStart int64    `json:"seconds_start"`
	SecondsEnd   int64    `json:"seconds_end"`
}

// Filter ???
type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
