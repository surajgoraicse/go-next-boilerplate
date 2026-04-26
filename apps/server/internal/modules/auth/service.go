package auth

import (
	"context"

	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
	db_sqlc "github.com/surajgoraicse/go-next-boilerplate/internal/db/sqlc"
)

type Service struct {
	queries *db_sqlc.Querier
	config  *config.Config
}

func NewService(queries db_sqlc.Querier, config *config.Config) *Service {
	return &Service{
		queries: &queries,
		config:  config,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (int, error) {

	return 0, nil
}
