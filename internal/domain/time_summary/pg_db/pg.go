package pg_db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// PgTimeSummaryRepository ???
type PgTimeSummaryRepository struct {
	db *sql.DB
}

// NewPgRepository ???
func NewPgRepository(db *sql.DB) *PgTimeSummaryRepository {
	return &PgTimeSummaryRepository{db: db}
}

// CreateTimeSummary ???
func (r *PgTimeSummaryRepository) CreateTimeSummary(
	ctx context.Context,
	ts *time_summary.TimeSummary,
) error {
	breaksJSON, err := json.Marshal(ts.Breaks)
	if err != nil {
		return fmt.Errorf("could not marshal breaks to json: %w", err)
	}

	builderInsert := sq.Insert("time_summary").
		PlaceholderFormat(sq.Dollar).
		Options("IGNORE").
		Columns("mac_address", "date", "seconds", "breaks", "seconds_begin", "seconds_end").
		Values(ts.MacAddress,
			ts.Date,
			ts.Seconds,
			breaksJSON,
			ts.SecondsStart,
			ts.SecondsEnd,
		)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return fmt.Errorf("create time query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert time: %w", err)
	}

	return nil
}

// GetTimeSummary ???
func (r *PgTimeSummaryRepository) GetTimeSummary(
	ctx context.Context,
	filters []*time_summary.Filter,
) ([]*time_summary.TimeSummary, error) {
	cond := make([]string, 0, len(filters))
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
		var breaksJSON []byte

		var ts time_summary.TimeSummary
		err = rows.Scan(
			&ts.MacAddress,
			&ts.Seconds,
			&breaksJSON,
			&ts.Date,
			&ts.SecondsStart,
			&ts.SecondsEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan time summary: %v", err)
		}

		err = json.Unmarshal(breaksJSON, &ts.Breaks)
		if err != nil {
			return nil, fmt.Errorf("error unmarshal breaks to json: %v", err)
		}

		res = append(res, &ts)
	}

	return res, nil
}
