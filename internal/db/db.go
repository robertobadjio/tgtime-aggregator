package db

import (
	"database/sql"
	"time"

	// Register some standard stuff
	_ "github.com/lib/pq"
	"github.com/robertobadjio/tgtime-aggregator/internal/config"
)

// GetDB ???
func GetDB() *sql.DB {
	pgCfg, err := config.NewPGConfig()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", pgCfg.DSN())
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}
