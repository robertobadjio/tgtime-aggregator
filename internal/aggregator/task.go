package aggregator

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/robertobadjio/tgtime-aggregator/internal/db"
	implementationT "github.com/robertobadjio/tgtime-aggregator/internal/domain/time/implementation"
	tPgRepo "github.com/robertobadjio/tgtime-aggregator/internal/domain/time/pg_db"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
	implementationTs "github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary/implementation"
	tsPgRepo "github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary/pg_db"
)

// Aggregate ???
// TODO: Горизонтальное масштабирование
func Aggregate() {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 0, 1, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	for {
		time.Sleep(d)
		d = 24 * time.Hour

		_ = calcTimeSummary() // TODO: Handle error
	}
}

func calcTimeSummary() error {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	tService := implementationT.NewTimeService(tPgRepo.NewPgRepository(db.GetDB()), logger)
	tsService := implementationTs.NewTimeSummaryService(tsPgRepo.NewPgRepository(db.GetDB()), logger)

	agr := NewAggregator(getPreviousDate("Europe/Moscow"), tService)
	ctx := context.TODO()
	macAddresses, err := tService.GetMacAddresses(ctx, getPreviousDate("Europe/Moscow"))
	if err != nil {
		return fmt.Errorf("error calc time summary: %w", err)
	}
	for _, macAddress := range macAddresses {
		timeSummary, err := calcTimeSummaryByMacAddress(ctx, agr, macAddress)
		if err != nil {
			_ = logger.Log("msg", err.Error())
		}

		err = tsService.CreateTimeSummary(ctx, timeSummary)
		if err != nil {
			return fmt.Errorf("error calc time summary by mac address "+macAddress+": %w", err)
		}
	}

	return nil
}

func calcTimeSummaryByMacAddress(
	ctx context.Context,
	agr *Aggregator,
	macAddress string,
) (time_summary.TimeSummary, error) {
	timeSummary, err := agr.AggregateTime(ctx, macAddress)
	if err != nil {
		return time_summary.TimeSummary{}, fmt.Errorf("error calc time summary by mac address "+macAddress+": %w", err)

	}

	return timeSummary, nil
}

func getPreviousDate(location string) time.Time {
	moscowLocation, _ := time.LoadLocation(location)
	return time.Now().AddDate(0, 0, -1).In(moscowLocation)
}
