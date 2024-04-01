package time

import (
	"context"
)

type Service interface {
	CreateTime(ctx context.Context, t *TimeUser) error
}
