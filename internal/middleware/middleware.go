package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type MiddlewareFunc func(xCtx *fiber.Ctx) appctx.Response

// FilterFunc is a iterator resolver in each middleware registered
func FilterFunc(conf *config.Config, xCtx *fiber.Ctx, mfs []MiddlewareFunc) appctx.Response {
	// Initiate postive case
	var response = appctx.Response{Code: fiber.StatusOK}
	for _, mf := range mfs {
		if response = mf(xCtx); response.Code != fiber.StatusOK {
			return response
		}
	}

	return response
}
