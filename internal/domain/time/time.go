package time

// Time ???
type Time struct {
	MacAddress string `json:"mac_address"`
	Seconds    int64  `json:"seconds"`
	RouterID   int64  `json:"router_id"`
}
