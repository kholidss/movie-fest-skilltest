package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
)

type userAuth struct {
	cfg *config.Config
}

func NewUUserAuthMiddleware(cfg *config.Config) *userAuth {
	return &userAuth{
		cfg: cfg,
	}
}

func (u *userAuth) Authenticate(xCtx *fiber.Ctx) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx)
		lf        = logger.NewFields(
			logger.EventName("UserAuthMiddleware"),
			logger.Any("X-Request-ID", requestID),
		)
	)

	ctx, span := tracer.NewSpan(xCtx.Context(), "middleware.user_auth", nil)
	defer span.End()

	logger.InfoWithContext(ctx, "success authenticate user", lf...)
	return *appctx.NewResponse().WithCode(fiber.StatusOK)
}
