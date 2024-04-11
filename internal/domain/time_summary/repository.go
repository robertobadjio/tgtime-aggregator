package time_summary

import (
	"context"
)

type Repository interface {
	CreateTimeSummary(ctx context.Context, timeSummary *TimeSummary) error
	GetTimeSummaryByDate(ctx context.Context, macAddress string, date string) (*TimeSummary, error)
	GetTimeSummaryAllByDate(ctx context.Context, date string) ([]*TimeSummary, error)
}
