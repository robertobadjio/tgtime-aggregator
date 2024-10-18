package pg_db

import (
	"database/sql"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// PgTimeSummaryRepository ???
type PgTimeSummaryRepository struct {
	db *sql.DB
}

// NewPgRepository Конструктор PostgresQL репозитория.
func NewPgRepository(db *sql.DB) time_summary.Repository {
	return &PgTimeSummaryRepository{db: db}
}
