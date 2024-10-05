package time

import (
	"context"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

type Service interface {
	CreateTime(ctx context.Context, time *time.TimeUser) (*time.TimeUser, error)
	GetTimeSummary(ctx context.Context, filters []*time_summary.Filter) ([]*time_summary.TimeSummary, error)
	ServiceStatus(ctx context.Context) int
}
