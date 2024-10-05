package kafka

type PreviousDayInfoMessage struct {
	MacAddress   string `json:"mac_address"`
	Seconds      int64  `json:"seconds"`
	BreaksJson   []byte `json:"breaks"`
	Date         string `json:"date"`
	SecondsStart int64  `json:"seconds_start"`
	SecondsEnd   int64  `json:"seconds_end"`
}

const PreviousDayInfoTopic = "previous-day-info"
const partition = 0
