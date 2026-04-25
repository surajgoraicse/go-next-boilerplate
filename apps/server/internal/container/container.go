package container

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/logger"
	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
	"github.com/surajgoraicse/go-next-boilerplate/internal/db"
)

type Container struct {
	// system
	Config *config.Config
	Logger logger.Logger

	// DB
	DB *pgxpool.Pool

	// ---------- modules ----------
	// auth

}

func NewContainer(ctx context.Context) *Container {
	// config setup
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load the config: " + err.Error())
	}

	// logger setup
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		panic("failed to initialize logger : error : " + err.Error())
	}

	// db setup
	dbService := db.NewDatabaseService(cfg)
	db, err := dbService.Connect(ctx)
	if err != nil {
		panic("failed to connect to db : error : " + err.Error())
	}
	logger.Info("db connected successfully")

	logger.Info("all services initialized successfully")
	return &Container{
		Config: cfg,
		Logger: logger,
		DB:     db,
	}
}

// close all the resources
func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
	func() {
		_ = c.Logger.Sync()
	}()
}
