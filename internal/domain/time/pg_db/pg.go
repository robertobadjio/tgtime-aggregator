package pg_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	timeDomain "github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
)

// PgTimeRepository ???
type PgTimeRepository struct {
	db *sql.DB
}

// NewPgRepository ???
func NewPgRepository(db *sql.DB) *PgTimeRepository {
	return &PgTimeRepository{db: db}
}

// CreateTime ???
func (r *PgTimeRepository) CreateTime(ctx context.Context, t *timeDomain.Time) error {
	builderInsert := sq.Insert("time").
		PlaceholderFormat(sq.Dollar).
		Columns("mac_address", "seconds", "router_id").
		Values(t.MacAddress, t.Seconds, t.RouterID)

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

// GetByFilters ???
func (r *PgTimeRepository) GetByFilters(ctx context.Context, query timeDomain.Query) ([]*timeDomain.Time, error) {
	cond := make([]string, 0, 4)
	if query.SecondsStart != 0 {
		cond = append(cond, fmt.Sprintf("seconds::integer >= %d", query.SecondsStart))
	}
	if query.SecondsEnd != 0 {
		cond = append(cond, fmt.Sprintf("seconds::integer <= %d", query.SecondsEnd))
	}
	if query.MacAddress != "" {
		cond = append(cond, fmt.Sprintf(`mac_address = '%s'`, query.MacAddress))
	}
	if query.RouterID != 0 {
		cond = append(cond, fmt.Sprintf(`router_id = %d`, query.RouterID))
	}
	rows, err := r.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT mac_address, seconds FROM time WHERE %s ORDER BY seconds",
			strings.Join(cond, " AND "),
		),
	)
	if err != nil {
		return []*timeDomain.Time{}, nil
	}
	defer func() {
		_ = rows.Close()
	}()

	times := make([]*timeDomain.Time, 0)
	for rows.Next() {
		t := new(timeDomain.Time)
		err = rows.Scan(&t.MacAddress, &t.Seconds)
		if err != nil {
			return []*timeDomain.Time{}, nil
		}

		times = append(times, t)
	}

	return times, nil
}

// GetSecondsDayByDate ???
func (r *PgTimeRepository) GetSecondsDayByDate(ctx context.Context, query timeDomain.Query, sort string) (int64, error) {
	var beginSecond int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT seconds FROM time WHERE mac_address = $1 AND seconds::integer BETWEEN $2 AND $3 ORDER BY seconds "+sort+" LIMIT 1",
		query.MacAddress,
		query.SecondsStart,
		query.SecondsEnd,
	).Scan(&beginSecond)
	if err == nil {
		return beginSecond, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	return 0, fmt.Errorf("error getting time summary from db: %v", err)
}

// GetMacAddresses ???
func (r *PgTimeRepository) GetMacAddresses(ctx context.Context, query timeDomain.Query) ([]string, error) {
	cond := make([]string, 0, 4)
	if query.SecondsStart != 0 {
		cond = append(cond, fmt.Sprintf("seconds::integer >= %d", query.SecondsStart))
	}
	if query.SecondsEnd != 0 {
		cond = append(cond, fmt.Sprintf("seconds::integer <= %d", query.SecondsEnd))
	}
	if query.MacAddress != "" {
		cond = append(cond, fmt.Sprintf(`mac_address = "%s"`, query.MacAddress))
	}
	if query.RouterID != 0 {
		cond = append(cond, fmt.Sprintf(`router_id = "%d"`, query.RouterID))
	}

	rows, err := r.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT mac_address FROM time WHERE %s GROUP BY mac_address",
			strings.Join(cond, " AND "),
		),
	)
	if err != nil {
		return []string{}, fmt.Errorf("error getting mac addresses from db: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	res := make([]string, 0)
	for rows.Next() {
		var macAddress string
		err = rows.Scan(&macAddress)
		if err != nil {
			return []string{}, nil
		}

		res = append(res, macAddress)
	}

	return res, nil
}
