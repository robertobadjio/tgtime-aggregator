package time_summary

import (
	"context"
)

// Repository ???
type Repository interface {
	CreateTimeSummary(ctx context.Context, timeSummary TimeSummary) error
	GetByFilters(ctx context.Context, filters []Filter) ([]TimeSummary, error)
}
