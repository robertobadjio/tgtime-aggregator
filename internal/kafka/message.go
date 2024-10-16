package kafka

// Break ???
type Break struct {
	SecondsStart int64
	SecondsEnd   int64
}

// PreviousDayInfoMessage ???
type PreviousDayInfoMessage struct {
	MacAddress   string   `json:"mac_address"`
	Seconds      int64    `json:"seconds"`
	Breaks       []*Break `json:"breaks"`
	Date         string   `json:"date"`
	SecondsStart int64    `json:"seconds_start"`
	SecondsEnd   int64    `json:"seconds_end"`
}

// const PreviousDayInfoTopic = "previous-day-info"
const partition = 0
