package time

import (
	"context"
)

type Repository interface {
	CreateTime(ctx context.Context, t *TimeUser) error
}
