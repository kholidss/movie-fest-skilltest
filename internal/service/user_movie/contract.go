package usermovie

import (
	"context"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
)

type UserMovieService interface {
	Vote(ctx context.Context, authData presentation.UserAuthData, payload presentation.ReqUserMovieVote) appctx.Response
	UnVote(ctx context.Context, authData presentation.UserAuthData, payload presentation.ReqUserMovieVote) appctx.Response
}
