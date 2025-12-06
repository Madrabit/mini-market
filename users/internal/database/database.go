package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/madrabit/mini-market/users/internal/common"
	"time"
)

func ConnectDb() *sqlx.DB {
	cfg, err := common.Load()
	if err != nil {
		panic(err)
	}
	return ConnectDbWithCfg(cfg)
}

func ConnectDbWithCfg(cfg common.Config) *sqlx.DB {
	db := sqlx.MustConnect(cfg.DB.Database, cfg.DB.DSN())
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db
}
