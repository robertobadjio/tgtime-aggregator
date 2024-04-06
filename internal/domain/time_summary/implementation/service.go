package implementation

import (
	"context"
	"github.com/go-kit/kit/log"
	"tgtime-aggregator/internal/domain/time_summary"
)

type TimeSummaryService struct {
	repository time_summary.Repository
	logger     log.Logger
}

func NewTimeSummaryService(
	rep time_summary.Repository,
	logger log.Logger,
) *TimeSummaryService {
	return &TimeSummaryService{
		repository: rep,
		logger:     logger,
	}
}

func (s *TimeSummaryService) CreateTimeSummary(
	ctx context.Context,
	ts *time_summary.TimeSummary,
) error {
	if err := s.repository.CreateTimeSummary(ctx, ts); err != nil {
		s.logger.Log("msg", err.Error())
		return err // TODO: !
	}

	return nil
}
