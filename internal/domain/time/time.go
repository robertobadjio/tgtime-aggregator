package time

// TimeUser ???
type TimeUser struct {
	MacAddress string `json:"mac_address"`
	Seconds    int64  `json:"seconds"`
	RouterID   int64  `json:"router_id"`
}
