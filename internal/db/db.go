package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"tgtime-aggregator/internal/config"
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

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}
