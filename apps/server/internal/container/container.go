package container

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/logger"
	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
	"github.com/surajgoraicse/go-next-boilerplate/internal/db"
	db_sqlc "github.com/surajgoraicse/go-next-boilerplate/internal/db/sqlc"
	"github.com/surajgoraicse/go-next-boilerplate/internal/modules/auth"
)

type Container struct {
	// system
	Config *config.Config
	Logger logger.Logger

	// DB
	DB *pgxpool.Pool

	// ----------  modules ----------
	// auth
	AuthService *auth.Service
	AuthHandler *auth.Handler
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

	queries := db_sqlc.New(db)

	// ------ modules initialization ------
	// auth module
	authService := auth.NewService(queries, cfg)
	authHandler := auth.NewHandler(authService)

	logger.Info("all services initialized successfully")
	return &Container{
		Config:      cfg,
		Logger:      logger,
		DB:          db,
		AuthHandler: authHandler,
		AuthService: authService,
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
