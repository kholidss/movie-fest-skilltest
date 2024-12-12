package usermovie

import (
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type userMovieService struct {
	cfg               *config.Config
	repoMovie         repositories.MovieRepository
	repoMovieVote     repositories.MovieVoteRepository
	repoActionHistory repositories.ActionHistoryRepository
}

func NewSvcUserMovie(
	cfg *config.Config,
	repoMovie repositories.MovieRepository,
	repoMovieVote repositories.MovieVoteRepository,
	repoActionHistory repositories.ActionHistoryRepository,
) UserMovieService {
	return &userMovieService{
		cfg:               cfg,
		repoMovie:         repoMovie,
		repoMovieVote:     repoMovieVote,
		repoActionHistory: repoActionHistory,
	}
}
