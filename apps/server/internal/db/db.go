package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
)

type DatabaseService struct {
	config *config.Config
}

func NewDatabaseService(cfg *config.Config) *DatabaseService {
	return &DatabaseService{
		config: cfg,
	}
}

type dbConnectionParams struct {
	DBHost                string
	DBPort                string
	DBUser                string
	DBPassword            string
	DBName                string
	SSLMode               string
	DBMaxConn             int32
	DBMinConn             int32
	DBConnMaxLifetime     time.Duration
	DBConnMaxIdleLifetime time.Duration
	DBHealthCheckPeriod   time.Duration
	ConnectTimeout        time.Duration
}

func (ds *DatabaseService) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	dbConfig, err := ds.withPgxConfig()
	if err != nil {
		return nil, fmt.Errorf("error parsing database configuration: %v", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}

func (ds *DatabaseService) withPgxConfig() (*pgxpool.Config, error) {
	dbURLConfig, err := ds.loadDBConnectionConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading database connection config: %v", err)
	}

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbURLConfig.DBHost, dbURLConfig.DBPort, dbURLConfig.DBUser, dbURLConfig.DBPassword, dbURLConfig.DBName, dbURLConfig.SSLMode)
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing database URL: %v", err)
	}

	dbConfig.MaxConns = dbURLConfig.DBMaxConn
	dbConfig.MinConns = dbURLConfig.DBMinConn
	dbConfig.MaxConnLifetime = dbURLConfig.DBConnMaxLifetime
	dbConfig.MaxConnIdleTime = dbURLConfig.DBConnMaxIdleLifetime
	dbConfig.HealthCheckPeriod = dbURLConfig.DBHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = dbURLConfig.ConnectTimeout

	// Ensure all connections use swags_me schema by default
	if dbConfig.ConnConfig.RuntimeParams == nil {
		dbConfig.ConnConfig.RuntimeParams = make(map[string]string)
	}
	dbConfig.ConnConfig.RuntimeParams["search_path"] = "swags_me,public"

	return dbConfig, nil
}

func (ds *DatabaseService) loadDBConnectionConfig() (*dbConnectionParams, error) {
	return &dbConnectionParams{
		DBHost:                ds.config.DBHost,
		DBPort:                ds.config.DBPort,
		DBUser:                ds.config.DBUser,
		DBPassword:            ds.config.DBPassword,
		DBName:                ds.config.DBName,
		SSLMode:               ds.config.SSLMode,
		DBMaxConn:             ds.config.DBMaxConn,
		DBMinConn:             ds.config.DBMinConn,
		DBConnMaxLifetime:     ds.config.DBConnMaxLifetime,
		DBConnMaxIdleLifetime: ds.config.DBConnMaxIdleLifetime,
		DBHealthCheckPeriod:   ds.config.DBHealthCheckPeriod,
		ConnectTimeout:        ds.config.ConnectTimeout,
	}, nil
}
