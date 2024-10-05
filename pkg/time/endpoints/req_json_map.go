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

type GetTimeSummaryRequest struct {
	Filters []*time_summary.Filter `json:"filters"`
}
type GetTimeSummaryResponse struct {
	TimeSummary []*time_summary.TimeSummary `json:"time_summary"`
}

type ServiceStatusRequest struct{}
type ServiceStatusResponse struct {
	Code int `json:"status"`
}
