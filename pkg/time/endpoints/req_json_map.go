package endpoints

import (
	"tgtime-aggregator/internal/time"
)

type CreateTimeRequest struct {
	Time *time.TimeUser `json:"time"`
}
type CreateTimeResponse struct {
	Time *time.TimeUser `json:"time"`
}

type ServiceStatusRequest struct{}
type ServiceStatusResponse struct {
	Code int `json:"status"`
}
