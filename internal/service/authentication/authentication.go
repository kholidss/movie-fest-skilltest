package authentication

import (
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type authenticationService struct {
	cfg      *config.Config
	repoUser repositories.UserRepository
}

func NewSvcAuthentication(
	cfg *config.Config,
	repoUser repositories.UserRepository,
) AuthenticationService {
	return &authenticationService{
		cfg:      cfg,
		repoUser: repoUser,
	}
}
