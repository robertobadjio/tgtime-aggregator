package time

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/robertobadjio/tgtime-aggregator/internal/db"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time/implementation"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time/pg_db"
	"net/http"
	"os"
)

type apiService struct {
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

func NewService() Service {
	return &apiService{}
}

func (s *apiService) CreateTime(ctx context.Context, t *time.TimeUser) (*time.TimeUser, error) {
	repo := pg_db.NewPgRepository(db.GetDB())
	timeService := implementation.NewTimeService(repo, logger)
	err := timeService.CreateTime(ctx, t)
	if err != nil {
		logger.Log("msg", err.Error())
		return nil, fmt.Errorf("error saving time")
	}

	return t, nil
}

func (s *apiService) ServiceStatus(_ context.Context) int {
	logger.Log("msg", "Checking the Service health...")
	return http.StatusOK
}
