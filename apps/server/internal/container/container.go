package container

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/logger"
	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
	"go.uber.org/zap"
)

type Container struct {
	// system
	Config *config.Config
	Logger *zap.Logger

	// DB
	DB *pgxpool.Pool

	// ---------- modules ----------
	// auth

}

func NewContainer() *Container {
	// config setup
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load the config")
	}

	// logger setup
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		panic("failed to initialize logger : error : " + err.Error())
	}
	defer func() { _ = logger.Sync() }()

	// db setup

	return &Container{
		Config: cfg,
		Logger: logger,
	}
}
