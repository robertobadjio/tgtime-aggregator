package endpoints

import (
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

type CreateTimeRequest struct {
	Time *time.TimeUser `json:"time"`
}
type CreateTimeResponse struct {
	Time *time.TimeUser `json:"time"`
}

type GetTimeSummaryByDateRequest struct {
	MacAddress string `json:"mac_address"`
	Date       string `json:"date"`
}
type GetTimeSummaryByDateResponse struct {
	TimeSummary *time_summary.TimeSummary `json:"time_summary"`
}

type ServiceStatusRequest struct{}
type ServiceStatusResponse struct {
	Code int `json:"status"`
}
