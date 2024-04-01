package time_summary

import (
	"context"
)

type Repository interface {
	CreateTimeSummary(ctx context.Context, timeSummary *TimeSummary) error
}
