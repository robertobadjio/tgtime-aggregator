package implementation

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
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

func (s *TimeSummaryService) GetTimeSummaryByDate(
	ctx context.Context,
	macAddress, date string,
) (*time_summary.TimeSummary, error) {
	ts, err := s.repository.GetTimeSummaryByDate(ctx, macAddress, date)
	if err != nil {
		s.logger.Log("msg", err.Error())
		return nil, err
	}

	return ts, nil
}

func (s *TimeSummaryService) GetTimeSummaryAllByDate(
	ctx context.Context,
	date string,
) ([]*time_summary.TimeSummary, error) {
	ts, err := s.repository.GetTimeSummaryAllByDate(ctx, date)
	if err != nil {
		s.logger.Log("msg", err.Error())
		return nil, err
	}

	return ts, nil
}
