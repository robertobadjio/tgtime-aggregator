package aggregator

import (
	"context"
	"encoding/json"
	implementationT "github.com/robertobadjio/tgtime-aggregator/internal/domain/time/implementation"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
	"github.com/robertobadjio/tgtime-aggregator/internal/tgtime_api_client"
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
	user tgtime_api_client.User,
) (*time_summary.TimeSummary, error) {
	times, err := agg.timeService.GetByFilters(ctx, user.MacAddress, agg.date, 0)
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
	begin, err := agg.timeService.GetStartSecondDayByDate(ctx, user.MacAddress, agg.date)
	if err != nil {
		return nil, err
	}
	end, err := agg.timeService.GetEndSecondDayByDate(ctx, user.MacAddress, agg.date)
	if err != nil {
		return nil, err
	}

	return &time_summary.TimeSummary{
		MacAddress:   user.MacAddress,
		Date:         agg.date,
		Seconds:      seconds,
		BreaksJson:   breaksJson,
		SecondsStart: begin,
		SecondsEnd:   end,
	}, nil
}
