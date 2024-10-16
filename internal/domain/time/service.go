package time

import (
	"context"
	"time"

	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// Service ???
type Service interface {
	CreateTime(ctx context.Context, t *TimeUser) error
	GetByFilters(
		ctx context.Context,
		macAddress string,
		date time.Time,
		routerID int,
	) ([]*TimeUser, error)
	GetStartSecondDayByDate(
		ctx context.Context,
		macAddress string,
		date time.Time,
	) (int64, error)
	AggregateDayTotalTime(times []*TimeUser) (int64, error)
	GetAllBreaksByTimesOld(times []*TimeUser) ([]*time_summary.Break, error)
	GetMacAddresses(ctx context.Context) ([]string, error)
}
