package time

import (
	"context"
	"time"
)

type Service interface {
	CreateTime(ctx context.Context, t *TimeUser) error
	GetByFilters(
		ctx context.Context,
		macAddress string,
		date time.Time,
		routerId int,
	) ([]*TimeUser, error)
	GetStartSecondDayByDate(
		ctx context.Context,
		macAddress string,
		date time.Time,
	) (int64, error)
	AggregateDayTotalTime(times []*TimeUser) (int64, error)
	GetAllBreaksByTimesOld(times []*TimeUser) ([]*Break, error)
	GetMacAddresses(ctx context.Context) ([]string, error)
}
