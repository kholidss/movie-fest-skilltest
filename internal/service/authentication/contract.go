package authentication

import (
	"context"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
)

type AuthenticationService interface {
	RegisterUser(ctx context.Context, payload presentation.ReqRegisterUser) appctx.Response
	LoginUser(ctx context.Context, payload presentation.ReqLoginUser) appctx.Response
}
