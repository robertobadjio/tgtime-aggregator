package time

import (
	"context"
	"time"

	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// Service ???
type Service interface {
	CreateTime(ctx context.Context, t *Time) error
	GetByFilters(
		ctx context.Context,
		macAddress string,
		date time.Time,
		routerID int,
	) ([]*Time, error)
	GetStartSecondDayByDate(
		ctx context.Context,
		macAddress string,
		date time.Time,
	) (int64, error)
	AggregateDayTotalTime(times []*Time) (int64, error)
	GetAllBreaksByTimesOld(times []*Time) ([]*time_summary.Break, error)
	GetMacAddresses(ctx context.Context) ([]string, error)
}
