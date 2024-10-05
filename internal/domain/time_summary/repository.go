package time_summary

import (
	"context"
)

type Repository interface {
	CreateTimeSummary(ctx context.Context, timeSummary *TimeSummary) error
	GetTimeSummary(ctx context.Context, filters []*Filter) ([]*TimeSummary, error)
}
