package time

import (
	"context"
	"fmt"
	"net/http"
	"os"
	t "time"

	"github.com/go-kit/kit/log"
	aggregator2 "github.com/robertobadjio/tgtime-aggregator/internal/aggregator"
	"github.com/robertobadjio/tgtime-aggregator/internal/db"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time/implementation"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time/pg_db"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
	timeSummaryimplementation "github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary/implementation"
	domainTimeSummary "github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary/pg_db"
)

type apiService struct {
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
}

// NewService ???
func NewService() Service {
	return &apiService{}
}

// CreateTime ???
func (s *apiService) CreateTime(ctx context.Context, t *time.Time) (*time.Time, error) {
	repo := pg_db.NewPgRepository(db.GetDB())
	timeService := implementation.NewTimeService(repo, logger)
	err := timeService.CreateTime(ctx, t)
	if err != nil {
		_ = logger.Log("msg", err.Error())
		return nil, fmt.Errorf("error saving time")
	}

	return t, nil
}

// GetTimeSummary ???
func (s *apiService) GetTimeSummary(
	ctx context.Context,
	filters []time_summary.Filter,
) ([]time_summary.TimeSummary, error) {
	repo := domainTimeSummary.NewPgRepository(db.GetDB())
	timeSummaryService := timeSummaryimplementation.NewTimeSummaryService(repo, logger)
	ts, err := timeSummaryService.GetByFilters(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error getting time summary")
	}

	// TODO: Refactoring
	flagMacAddress := ""
	flagToday := false
	for _, filter := range filters {
		if filter.Key == "mac_address" {
			flagMacAddress = filter.Value
		} else if filter.Key == "date" && filter.Value == getDate("Europe/Moscow").Format("2006-01-02") {
			flagToday = true
		}
	}
	if flagMacAddress != "" && flagToday {
		tService := implementation.NewTimeService(pg_db.NewPgRepository(db.GetDB()), logger)
		agr := aggregator2.NewAggregator(getDate("Europe/Moscow"), tService)
		todayTimeSummary, _ := agr.AggregateTime(ctx, flagMacAddress) // TODO: Handle error
		ts = append(ts, todayTimeSummary)
	}

	return ts, nil
}

func getDate(location string) t.Time {
	moscowLocation, _ := t.LoadLocation(location)
	return t.Now().In(moscowLocation)
}

// ServiceStatus ???
func (s *apiService) ServiceStatus(_ context.Context) int {
	_ = logger.Log("msg", "Checking the Service health...")
	return http.StatusOK
}
