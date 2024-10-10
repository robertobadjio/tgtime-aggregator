package pg_db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
	"strings"
)

type PgTimeSummaryRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgTimeSummaryRepository {
	return &PgTimeSummaryRepository{db: db}
}

func (r *PgTimeSummaryRepository) CreateTimeSummary(ctx context.Context, ts *time_summary.TimeSummary) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT IGNORE INTO time_summary (mac_address, date, seconds, breaks, seconds_begin, seconds_end) VALUES ($1, $2, $3, $4, $5, $6)",
		ts.MacAddress,
		ts.Date,
		ts.Seconds,
		ts.BreaksJson,
		ts.SecondsStart,
		ts.SecondsEnd,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgTimeSummaryRepository) GetTimeSummary(
	ctx context.Context,
	filters []*time_summary.Filter,
) ([]*time_summary.TimeSummary, error) {
	cond := make([]string, 0, len(filters)) // TODO: !
	for _, filter := range filters {
		cond = append(cond, fmt.Sprintf("%s = '%s'", filter.Key, filter.Value))
	}

	rows, err := r.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT mac_address, seconds, breaks, date, seconds_begin, seconds_end FROM time_summary WHERE %s",
			strings.Join(cond, " AND "),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("error getting time summary: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var res []*time_summary.TimeSummary
	for rows.Next() {
		var ts time_summary.TimeSummary
		err = rows.Scan(
			&ts.MacAddress,
			&ts.Seconds,
			&ts.BreaksJson,
			&ts.Date,
			&ts.SecondsStart,
			&ts.SecondsEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan time summary: %v", err)
		}

		res = append(res, &ts)
	}

	return res, nil
}
