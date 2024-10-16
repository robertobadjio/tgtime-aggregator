package time

import (
	"context"
)

// Query ???
type Query struct {
	MacAddress   string
	SecondsStart int64
	SecondsEnd   int64
	RouterID     int
}

// Repository ???
type Repository interface {
	CreateTime(ctx context.Context, t *Time) error
	GetByFilters(ctx context.Context, query Query) ([]*Time, error)
	GetSecondsDayByDate(ctx context.Context, query Query, sort string) (int64, error)
	GetMacAddresses(ctx context.Context, query Query) ([]string, error)
}
