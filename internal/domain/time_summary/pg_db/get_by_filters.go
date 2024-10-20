package pg_db

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// GetByFilters ???
func (r *PgTimeSummaryRepository) GetByFilters(
	ctx context.Context,
	filters []time_summary.Filter,
) ([]time_summary.TimeSummary, error) {
	// TODO: Пагинация
	builderSelect := sq.Select("mac_address", "seconds", "breaks", "date", "seconds_begin", "seconds_end").
		From("time_summary").
		PlaceholderFormat(sq.Dollar).
		OrderBy("date ASC")

	for _, filter := range filters {
		builderSelect = builderSelect.Where(sq.Eq{filter.Key: filter.Value})
	}

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("create get time summaries query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting time summary: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var res []time_summary.TimeSummary
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

		res = append(res, ts)
	}

	return res, nil
}
