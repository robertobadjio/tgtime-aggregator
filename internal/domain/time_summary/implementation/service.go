package implementation

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// TimeSummaryService ???
type timeSummaryService struct {
	repository time_summary.Repository
	logger     log.Logger
}

// NewTimeSummaryService ???
func NewTimeSummaryService(
	rep time_summary.Repository,
	logger log.Logger,
) time_summary.Service {
	return &timeSummaryService{
		repository: rep,
		logger:     logger,
	}
}

// CreateTimeSummary ???
func (s *timeSummaryService) CreateTimeSummary(
	ctx context.Context,
	ts time_summary.TimeSummary,
) error {
	if err := s.repository.CreateTimeSummary(ctx, ts); err != nil {
		_ = s.logger.Log("msg", err.Error())
		return err // TODO: !
	}

	return nil
}

// GetByFilters ???
func (s *timeSummaryService) GetByFilters(
	ctx context.Context,
	filters []time_summary.Filter,
) ([]time_summary.TimeSummary, error) {
	ts, err := s.repository.GetByFilters(ctx, filters)
	if err != nil {
		_ = s.logger.Log("msg", err.Error())
		return nil, err
	}

	return ts, nil
}
