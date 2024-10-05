package pg_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	time2 "github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"strings"
)

type PgTimeRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgTimeRepository {
	return &PgTimeRepository{db: db}
}

func (r *PgTimeRepository) CreateTime(ctx context.Context, t *time2.TimeUser) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO time (mac_address, seconds, router_id) VALUES ($1, $2, $3)",
		t.MacAddress,
		t.Seconds,
		t.RouterId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgTimeRepository) GetByFilters(ctx context.Context, query time2.Query) ([]*time2.TimeUser, error) {
	cond := make([]string, 1, 4)
	if query.SecondsStart != 0 {
		cond = append(cond, fmt.Sprintf("seconds >= %d", query.SecondsStart))
	}
	if query.SecondsEnd != 0 {
		cond = append(cond, fmt.Sprintf("seconds <= %d", query.SecondsEnd))
	}
	if query.MacAddress != "" {
		cond = append(cond, fmt.Sprintf(`mac_address = "%s"`, query.MacAddress))
	}
	if query.RouterId != 0 {
		cond = append(cond, fmt.Sprintf(`router_id = "%d"`, query.RouterId))
	}

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT t.mac_address, t.second FROM time t WHERE "+strings.Join(cond, " AND ")+" ORDER BY t.second",
	)
	if err != nil {
		return []*time2.TimeUser{}, nil
	}
	defer func() {
		_ = rows.Close()
	}()

	times := make([]*time2.TimeUser, 0)
	for rows.Next() {
		t := new(time2.TimeUser)
		err = rows.Scan(&t.MacAddress, &t.Seconds)
		if err != nil {
			return []*time2.TimeUser{}, nil
		}

		times = append(times, t)
	}

	return times, nil
}

func (r *PgTimeRepository) GetSecondsDayByDate(ctx context.Context, query time2.Query, sort string) (int64, error) {
	var beginSecond int64
	if err := r.db.QueryRowContext(
		ctx,
		"SELECT t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3 ORDER BY t.second "+sort+" LIMIT 1",
		query.MacAddress,
		query.SecondsStart,
		query.SecondsEnd,
	).Scan(&beginSecond); err == nil {
		return beginSecond, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else {
		return 0, fmt.Errorf("error getting time summary from db: %v", err)
	}
}

func (r *PgTimeRepository) GetMacAddresses(ctx context.Context, query time2.Query) ([]string, error) {
	cond := make([]string, 1, 4)
	if query.SecondsStart != 0 {
		cond = append(cond, fmt.Sprintf("seconds >= %d", query.SecondsStart))
	}
	if query.SecondsEnd != 0 {
		cond = append(cond, fmt.Sprintf("seconds <= %d", query.SecondsEnd))
	}
	if query.MacAddress != "" {
		cond = append(cond, fmt.Sprintf(`mac_address = "%s"`, query.MacAddress))
	}
	if query.RouterId != 0 {
		cond = append(cond, fmt.Sprintf(`router_id = "%d"`, query.RouterId))
	}

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT t.mac_address FROM time t WHERE "+strings.Join(cond, " AND ")+" GROUP BY t.mac_address",
	)
	if err != nil {
		return []string{}, nil
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
