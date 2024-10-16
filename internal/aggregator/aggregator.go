package aggregator

import (
	"context"
	"time"

	implementationT "github.com/robertobadjio/tgtime-aggregator/internal/domain/time/implementation"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// Aggregator Агрегатора времени сотрудника проведенного на работе / в офисе
type Aggregator struct {
	date        time.Time
	timeService *implementationT.TimeService
}

// NewAggregator Конструктор агрегатора времени
func NewAggregator(
	date time.Time,
	timeService *implementationT.TimeService,
) *Aggregator {
	return &Aggregator{
		date:        date,
		timeService: timeService,
	}
}

// AggregateTime Формирования итогового времени по сотруднику за предыдущий день
func (agg Aggregator) AggregateTime(
	ctx context.Context,
	macAddress string,
) (*time_summary.TimeSummary, error) {
	times, err := agg.timeService.GetByFilters(ctx, macAddress, agg.date, 0)
	if err != nil {
		return nil, err
	}
	seconds, err := agg.timeService.AggregateDayTotalTime(times)
	if err != nil {
		return nil, err
	}

	breaks, err := agg.timeService.GetAllBreaksByTimes(times)
	if err != nil {
		return nil, err
	}

	begin, err := agg.timeService.GetStartSecondDayByDate(ctx, macAddress, agg.date)
	if err != nil {
		return nil, err
	}
	end, err := agg.timeService.GetEndSecondDayByDate(ctx, macAddress, agg.date)
	if err != nil {
		return nil, err
	}

	return &time_summary.TimeSummary{
		MacAddress:   macAddress,
		Date:         agg.date.Format("2006-01-02"),
		Seconds:      seconds,
		Breaks:       breaks,
		SecondsStart: begin,
		SecondsEnd:   end,
	}, nil
}
