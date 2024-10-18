package endpoints

import (
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// CreateTimeRequest ???
type CreateTimeRequest struct {
	Time *time.Time `json:"time"`
}

// CreateTimeResponse ???
type CreateTimeResponse struct {
	Time *time.Time `json:"time"`
}

// GetTimeSummaryRequest ???
type GetTimeSummaryRequest struct {
	Filters []time_summary.Filter `json:"filters"`
}

// GetTimeSummaryResponse ???
type GetTimeSummaryResponse struct {
	TimeSummary []time_summary.TimeSummary `json:"time_summary"`
}

// ServiceStatusRequest ???
type ServiceStatusRequest struct{}

// ServiceStatusResponse ???
type ServiceStatusResponse struct {
	Code int `json:"status"`
}
