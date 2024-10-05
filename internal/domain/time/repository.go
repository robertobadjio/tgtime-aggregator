package time

import (
	"context"
)

type Query struct {
	MacAddress   string
	SecondsStart int64
	SecondsEnd   int64
	RouterId     int
}

type Repository interface {
	CreateTime(ctx context.Context, t *TimeUser) error
	GetByFilters(ctx context.Context, query Query) ([]*TimeUser, error)
	GetSecondsDayByDate(ctx context.Context, query Query, sort string) (int64, error)
	GetMacAddresses(ctx context.Context, query Query) ([]string, error)
}
