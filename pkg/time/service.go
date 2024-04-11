package time

import (
	"context"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

type Service interface {
	CreateTime(ctx context.Context, time *time.TimeUser) (*time.TimeUser, error)
	GetTimeSummaryByDate(ctx context.Context, macAddress string, date string) (*time_summary.TimeSummary, error)
	GetTimeSummaryAllByDate(ctx context.Context, date string) ([]*time_summary.TimeSummary, error)
	ServiceStatus(ctx context.Context) int
}
