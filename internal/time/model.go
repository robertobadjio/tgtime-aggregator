package time

type TimeUser struct {
	MacAddress string `json:"mac_address"`
	Seconds    int64  `json:"seconds"`
	RouterId   int8   `json:"router_id"`
}

type Break struct {
	BeginTime int64 `json:"beginTime"`
	EndTime   int64 `json:"endTime"`
}
