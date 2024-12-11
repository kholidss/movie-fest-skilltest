package cmsmovie

import (
	"context"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
)

type CMSMovieService interface {
	Create(ctx context.Context, authData presentation.UserAuthData, payload presentation.ReqCMSCreateMovie) appctx.Response
}
