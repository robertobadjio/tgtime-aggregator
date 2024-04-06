package time

import (
	"context"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
)

type Service interface {
	CreateTime(ctx context.Context, time *time.TimeUser) (*time.TimeUser, error)
	ServiceStatus(ctx context.Context) int
}
