package pg_db

import (
	"context"
	"database/sql"
	"tgtime-aggregator/internal/domain/time_summary"
)

type PgTimeSummaryRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgTimeSummaryRepository {
	return &PgTimeSummaryRepository{db: db}
}

func (r *PgTimeSummaryRepository) CreateTimeSummary(_ context.Context, ts *time_summary.TimeSummary) error {
	_, err := r.db.Exec(
		"INSERT INTO time_summary (mac_address, date, seconds, breaks, seconds_begin, seconds_end) VALUES ($1, $2, $3, $4, $5, $6)",
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
