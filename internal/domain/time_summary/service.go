package time_summary

import "context"

type Service interface {
	CreateTimeSummary(ctx context.Context, ts TimeSummary) error
}
