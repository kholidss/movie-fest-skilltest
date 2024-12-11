package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type httpHandlerFunc func(xCtx *fiber.Ctx, svc contract.Controller, conf *config.Config) appctx.Response

type Router interface {
	Route()
}
