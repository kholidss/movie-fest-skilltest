package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/bootstrap"
	"github.com/kholidss/movie-fest-skilltest/internal/controller"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/v1/authentication"
	"github.com/kholidss/movie-fest-skilltest/internal/handler"
	"github.com/kholidss/movie-fest-skilltest/internal/middleware"
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	moduleAuth "github.com/kholidss/movie-fest-skilltest/internal/service/authentication"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
)

type router struct {
	cfg   *config.Config
	fiber fiber.Router
}

func (rtr *router) handle(hfn httpHandlerFunc, svc contract.Controller, mdws ...middleware.MiddlewareFunc) fiber.Handler {
	return func(xCtx *fiber.Ctx) error {

		//check registered middleware functions
		if rm := middleware.FilterFunc(rtr.cfg, xCtx, mdws); rm.Code != fiber.StatusOK {
			// return response base on middleware
			res := *appctx.NewResponse().
				WithCode(rm.Code).
				WithError(rm.Errors).
				WithMessage(rm.Message)
			return rtr.response(xCtx, res)
		}

		//send to controller
		resp := hfn(xCtx, svc, rtr.cfg)
		return rtr.response(xCtx, resp)
	}
}

func (rtr *router) response(fiberCtx *fiber.Ctx, resp appctx.Response) error {
	fiberCtx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return fiberCtx.Status(resp.Code).Send(resp.Byte())
}

func (rtr *router) Route() {
	//init db
	db := bootstrap.RegistryMySQLDatabase(rtr.cfg)

	//define repositories
	repoUser := repositories.NewUserRepository(db)

	//define middleware
	//middlewareUserAuth := middleware.NewUserAuthMiddleware(rtr.cfg, repoUser)

	//define cdn
	_ = bootstrap.RegistryCDN(rtr.cfg)

	//define services
	svcRegisterUser := moduleAuth.NewSvcAuthentication(rtr.cfg, repoUser)

	//define controller
	ctrRegisterUser := authentication.NewRegisterUser(svcRegisterUser)

	health := controller.NewGetHealth()

	externalV1 := rtr.fiber.Group("/api/external/v1")

	pathAuthV1 := externalV1.Group("/auth")

	rtr.fiber.Get("/ping", rtr.handle(
		handler.HttpRequest,
		health,
	))

	//Path authentication
	pathAuthV1.Post("/register/user", rtr.handle(
		handler.HttpRequest,
		ctrRegisterUser,
	))

	//Path authentication

}

func NewRouter(cfg *config.Config, fiber fiber.Router) Router {
	return &router{
		cfg:   cfg,
		fiber: fiber,
	}
}
