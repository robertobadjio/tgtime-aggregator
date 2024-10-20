package pg_db

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// CreateTimeSummary ???
func (r *PgTimeSummaryRepository) CreateTimeSummary(
	ctx context.Context,
	ts time_summary.TimeSummary,
) error {
	breaksJSON, err := json.Marshal(ts.Breaks)
	if err != nil {
		return fmt.Errorf("could not marshal breaks to json: %w", err)
	}

	builderInsert := sq.Insert("time_summary").
		PlaceholderFormat(sq.Dollar).
		//Options("IGNORE").
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
