package aggregator

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"os"
	implementationT "tgtime-aggregator/internal/domain/time/implementation"
	tPgRepo "tgtime-aggregator/internal/domain/time/pg_db"
	"tgtime-aggregator/internal/domain/time_summary"
	implementationTs "tgtime-aggregator/internal/domain/time_summary/implementation"
	tsPgRepo "tgtime-aggregator/internal/domain/time_summary/pg_db"
	"tgtime-aggregator/internal/tgtime_api_client"
	"time"
)

var Db *sql.DB
var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

func AggregateTime() error {
	ctx := context.TODO()
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	date := time.Now().AddDate(0, 0, -1).In(moscowLocation)

	apiClient := tgtime_api_client.NewTgTimeClient()
	users, _ := apiClient.GetAllUsers()

	tRepo := tPgRepo.NewPgRepository(Db)
	tService := implementationT.NewTimeService(tRepo, logger)

	tsRepo := tsPgRepo.NewPgRepository(Db)
	tsService := implementationTs.NewTimeSummaryService(tsRepo, logger)

	for _, user := range users.Users {
		times, err := tService.GetByFilters(ctx, user.MacAddress, date, 0)
		seconds, err := tService.AggregateDayTotalTime(times)
		breaks, err := tService.GetAllBreaksByTimesOld(times)
		breaksJson, err := json.Marshal(breaks)
		begin, err := tService.GetStartSecondDayByDate(ctx, user.MacAddress, date)
		end, err := tService.GetEndSecondDayByDate(ctx, user.MacAddress, date)

		ts := time_summary.TimeSummary{
			MacAddress:   user.MacAddress,
			Date:         date,
			Seconds:      seconds,
			BreaksJson:   breaksJson,
			SecondsStart: begin,
			SecondsEnd:   end,
		}
		err = tsService.CreateTimeSummary(ctx, ts)
		if err != nil {
			return err
		}
	}

	return nil
}
