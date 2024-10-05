package aggregator

import (
	"context"
	"encoding/json"
	implementationT "github.com/robertobadjio/tgtime-aggregator/internal/domain/time/implementation"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
	"time"
)

type Aggregator struct {
	date        time.Time
	timeService *implementationT.TimeService
}

func NewAggregator(
	date time.Time,
	timeService *implementationT.TimeService,
) *Aggregator {
	return &Aggregator{
		date:        date,
		timeService: timeService,
	}
}

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
	breaks, err := agg.timeService.GetAllBreaksByTimesOld(times)
	if err != nil {
		return nil, err
	}
	breaksJson, err := json.Marshal(breaks)
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
		BreaksJson:   breaksJson,
		SecondsStart: begin,
		SecondsEnd:   end,
	}, nil
}
