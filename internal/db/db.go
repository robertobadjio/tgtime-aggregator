package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/robertobadjio/tgtime-aggregator/internal/config"
	"time"
)

func GetDB() *sql.DB {
	cfg := config.New()
	pgConString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DataBaseHost,
		cfg.DataBasePort,
		cfg.DataBaseUser,
		cfg.DataBasePassword,
		cfg.DataBaseName,
		cfg.DataBaseSslMode,
	)

	db, err := sql.Open("postgres", pgConString)
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
