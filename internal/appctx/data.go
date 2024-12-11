package appctx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type Data struct {
	FiberCtx *fiber.Ctx
	Cfg      *config.Config
}
