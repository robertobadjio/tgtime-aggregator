package implementation

import (
	"context"
	"github.com/go-kit/kit/log"
	"tgtime-aggregator/internal/time"
)

type timeService struct {
	repository time.Repository
	logger     log.Logger
}

func NewService(rep time.Repository, logger log.Logger) *timeService {
	return &timeService{
		repository: rep,
		logger:     logger,
	}
}

func (s *timeService) CreateTime(ctx context.Context, t *time.TimeUser) error {
	if err := s.repository.CreateTime(ctx, t); err != nil {
		s.logger.Log("msg", err.Error())
		return err // TODO: !
	}

	return nil
}
