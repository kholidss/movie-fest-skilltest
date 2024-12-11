package publicmovie

import (
	"context"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
)

type PublicMovieService interface {
	List(ctx context.Context, param presentation.ReqPublicMovieList) appctx.Response
}
