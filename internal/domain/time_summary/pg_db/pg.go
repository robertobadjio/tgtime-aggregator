package pg_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
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

func (r *PgTimeSummaryRepository) GetTimeSummaryByDate(
	_ context.Context,
	macAddress string,
	date string,
) (*time_summary.TimeSummary, error) {
	ts := new(time_summary.TimeSummary)

	if err := r.db.QueryRow(
		"SELECT mac_address, seconds, breaks, date, seconds_begin, seconds_end FROM time_summary WHERE mac_address = $1 AND date = $2",
		macAddress,
		date,
	).Scan(
		&ts.MacAddress,
		&ts.Seconds,
		&ts.BreaksJson,
		&ts.Date,
		&ts.SecondsStart,
		&ts.SecondsEnd,
	); err == nil {
		return ts, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else {
		return nil, fmt.Errorf("error getting time summary from db: %v", err)
	}
}

func (r *PgTimeSummaryRepository) GetTimeSummaryAllByDate(
	_ context.Context,
	date string,
) ([]*time_summary.TimeSummary, error) {
	rows, err := r.db.Query(
		"SELECT mac_address, seconds, breaks, date, seconds_begin, seconds_end FROM time_summary WHERE date = $1",
		date)
	if err != nil {
		return nil, fmt.Errorf("error getting time summary: %v", err)
	}
	defer rows.Close()

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
