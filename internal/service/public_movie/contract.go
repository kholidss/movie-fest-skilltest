package publicmovie

import (
	"context"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
)

type PublicMovieService interface {
	List(ctx context.Context, param presentation.ReqPublicMovieList) appctx.Response
	Track(ctx context.Context, payload presentation.ReqPublicTrackMovieViewer) appctx.Response
	Search(ctx context.Context, param presentation.ReqPublicMovieSearch) appctx.Response
}
