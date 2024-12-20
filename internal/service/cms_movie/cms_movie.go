package cmsmovie

import (
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/pkg/cdn"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type cmsMovieService struct {
	cfg               *config.Config
	repoMovie         repositories.MovieRepository
	repoGenre         repositories.GenreRepository
	repoMovieGenre    repositories.MovieGenreRepository
	repoActionHistory repositories.ActionHistoryRepository
	repoBucket        repositories.BucketRepository

	cdn cdn.CDN
}

func NewSvcCMSMovie(
	cfg *config.Config,
	repoMovie repositories.MovieRepository,
	repoGenre repositories.GenreRepository,
	repoMovieGenre repositories.MovieGenreRepository,
	repoActionHistory repositories.ActionHistoryRepository,
	repoBucket repositories.BucketRepository,
	cdn cdn.CDN,
) CMSMovieService {
	return &cmsMovieService{
		cfg:               cfg,
		repoMovie:         repoMovie,
		repoGenre:         repoGenre,
		repoMovieGenre:    repoMovieGenre,
		repoActionHistory: repoActionHistory,
		repoBucket:        repoBucket,
		cdn:               cdn,
	}
}
