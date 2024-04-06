package pg_db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"tgtime-aggregator/internal/db"
	time2 "tgtime-aggregator/internal/domain/time"
)

type PgTimeRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgTimeRepository {
	return &PgTimeRepository{db: db}
}

func (r *PgTimeRepository) CreateTime(_ context.Context, t *time2.TimeUser) error {
	_, err := db.GetDB().Exec(
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

func (r *PgTimeRepository) GetByFilters(_ context.Context, query time2.Query) ([]*time2.TimeUser, error) {
	var cond []string
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

	rows, err := db.GetDB().Query(
		"SELECT t.mac_address, t.second FROM time t WHERE " + strings.Join(cond, " AND ") + " ORDER BY t.second",
	)
	if err != nil {
		return []*time2.TimeUser{}, nil
	}
	defer rows.Close()

	times := make([]*time2.TimeUser, 0)
	for rows.Next() {
		t := new(time2.TimeUser)
		err = rows.Scan(&t.MacAddress, &t.Seconds)
		if err != nil {
			panic(err)
		}

		times = append(times, t)
	}

	return times, nil
}

func (r *PgTimeRepository) GetSecondsDayByDate(_ context.Context, query time2.Query, sort string) (int64, error) {
	var beginSecond int64
	row := db.GetDB().QueryRow(
		"SELECT t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3 ORDER BY t.second "+sort+" LIMIT 1",
		query.MacAddress,
		query.SecondsStart,
		query.SecondsEnd,
	)
	err := row.Scan(&beginSecond)

	if err == sql.ErrNoRows {
		return 0, err
	}

	return beginSecond, nil
}
