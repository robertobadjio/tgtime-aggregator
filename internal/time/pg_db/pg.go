package pg_db

import (
	"context"
	"database/sql"
	"tgtime-aggregator/internal/db"
	"tgtime-aggregator/internal/time"
)

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgRepository {
	return &PgRepository{db: db}
}

func (r *PgRepository) CreateTime(_ context.Context, t *time.TimeUser) error {
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
