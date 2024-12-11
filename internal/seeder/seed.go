package seeder

import (
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type seeder struct {
	cfg *config.Config
}

func NewSeedRun(cfg *config.Config) Seederer {
	return &seeder{
		cfg: cfg,
	}
}
