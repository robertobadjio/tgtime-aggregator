package time_summary

import "context"

// Service ???
type Service interface {
	CreateTimeSummary(ctx context.Context, ts TimeSummary) error
}
