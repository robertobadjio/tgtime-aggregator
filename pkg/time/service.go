package time

import (
	"context"
	"tgtime-aggregator/internal/time"
)

type Service interface {
	CreateTime(ctx context.Context, time *time.TimeUser) (*time.TimeUser, error)
	ServiceStatus(ctx context.Context) int
}
