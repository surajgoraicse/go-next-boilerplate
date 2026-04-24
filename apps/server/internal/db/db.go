package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
)

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	// configuration
	config, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		return nil, err
	}

	// set pool configuration
	if cfg.DBMaxConn > 0 {
		config.MaxConns = int32(cfg.DBMaxConn)
	}
	if cfg.DBMinConn > 0 {
		config.MinConns = int32(cfg.DBMinConn)
	}
	config.MaxConnLifetime = cfg.DBConnMaxLifetime
	config.MaxConnIdleTime = cfg.DBConnMaxIdleLifetime

	// create pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// verify connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	log.Printf("database connection established")
	return pool, nil
}
