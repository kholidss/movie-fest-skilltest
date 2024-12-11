package seeder

import (
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type seeder struct {
	cfg      *config.Config
	repoUser repositories.UserRepository
}

func NewSeedRun(cfg *config.Config, repoUser repositories.UserRepository) Seederer {
	return &seeder{
		cfg:      cfg,
		repoUser: repoUser,
	}
}
