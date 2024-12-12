package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/bootstrap"
	"github.com/kholidss/movie-fest-skilltest/internal/controller"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/v1/authentication"
	cmsmovie "github.com/kholidss/movie-fest-skilltest/internal/controller/v1/cms_movie"
	publicMovie "github.com/kholidss/movie-fest-skilltest/internal/controller/v1/public_movie"
	userMovie "github.com/kholidss/movie-fest-skilltest/internal/controller/v1/user_movie"
	"github.com/kholidss/movie-fest-skilltest/internal/handler"
	"github.com/kholidss/movie-fest-skilltest/internal/middleware"
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	moduleAuth "github.com/kholidss/movie-fest-skilltest/internal/service/authentication"
	moduleCMSMovie "github.com/kholidss/movie-fest-skilltest/internal/service/cms_movie"
	modulePublicMovie "github.com/kholidss/movie-fest-skilltest/internal/service/public_movie"
	moduleUserMovie "github.com/kholidss/movie-fest-skilltest/internal/service/user_movie"
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
	repoMovie := repositories.NewMovieRepository(db)
	repoMovieGenre := repositories.NewMovieGenreRepository(db)
	repoMovieVote := repositories.NewMovieVoteRepository(db)
	repoGenre := repositories.NewGenreRepository(db)
	repoActionHistory := repositories.NewActionHistoryRepository(db)
	repoBucket := repositories.NewBucketRepository(db)

	//define middleware
	middlewareAdminAuth := middleware.NewAdminAuthMiddleware(rtr.cfg, repoUser)
	middlewareUserAuth := middleware.NewUserAuthMiddleware(rtr.cfg, repoUser)

	//define cdn
	cdnStorage := bootstrap.RegistryCDN(rtr.cfg)

	//define services
	svcAuth := moduleAuth.NewSvcAuthentication(rtr.cfg, repoUser)
	svcCMSMovie := moduleCMSMovie.NewSvcCMSMovie(
		rtr.cfg,
		repoMovie,
		repoGenre,
		repoMovieGenre,
		repoActionHistory,
		repoBucket,
		cdnStorage,
	)
	svcPublicMovie := modulePublicMovie.NewSvcCMSMovie(
		rtr.cfg,
		repoMovie,
		repoGenre,
		repoMovieGenre,
		repoActionHistory,
		repoBucket,
	)

	svcUserMovie := moduleUserMovie.NewSvcUserMovie(
		rtr.cfg,
		repoMovie,
		repoMovieVote,
		repoActionHistory,
	)

	//define controller
	ctrHealthCheck := controller.NewGetHealth()
	ctrRegisterUser := authentication.NewRegisterUser(svcAuth)
	ctrLoginUser := authentication.NewLoginUser(svcAuth)
	ctrLoginAdmin := authentication.NewLoginAdmin(svcAuth)
	ctrCMSCreateMovie := cmsmovie.NewCMSMovieCreate(svcCMSMovie)
	ctrCMSUpdateMovie := cmsmovie.NewCMSMovieUpdate(svcCMSMovie)
	ctrCMSMostView := cmsmovie.NewCMSMostView(svcCMSMovie)
	ctrPublicMovieList := publicMovie.NewPublicListMovie(svcPublicMovie)
	ctrPublicTrackMovieViewer := publicMovie.NewPublicTrackMovieViewer(svcPublicMovie)
	ctrPublicMovieSearch := publicMovie.NewPublicMovieSearch(svcPublicMovie)
	ctrUserMovieVote := userMovie.NewUserVoteMovie(svcUserMovie)

	externalV1 := rtr.fiber.Group("/api/external/v1")
	pathAuthV1 := externalV1.Group("/auth")
	pathCMSMovieV1 := externalV1.Group("/cms/movie")
	pathCMSListV1 := externalV1.Group("/cms/list")
	pathPublicMovie := externalV1.Group("/public/movie")
	pathUserMovie := externalV1.Group("/user/movie")

	rtr.fiber.Get("/ping", rtr.handle(
		handler.HttpRequest,
		ctrHealthCheck,
	))

	//Path authentication
	pathAuthV1.Post("/register/user", rtr.handle(
		handler.HttpRequest,
		ctrRegisterUser,
	))
	pathAuthV1.Post("/login/user", rtr.handle(
		handler.HttpRequest,
		ctrLoginUser,
	))
	pathAuthV1.Post("/login/admin", rtr.handle(
		handler.HttpRequest,
		ctrLoginAdmin,
	))

	//Path cms movie
	pathCMSMovieV1.Post("/create", rtr.handle(
		handler.HttpRequest,
		ctrCMSCreateMovie,
		middlewareAdminAuth.Authenticate,
	))
	pathCMSMovieV1.Put("/update/:id", rtr.handle(
		handler.HttpRequest,
		ctrCMSUpdateMovie,
		middlewareAdminAuth.Authenticate,
	))

	//Path cms list
	pathCMSListV1.Get("/most-view", rtr.handle(
		handler.HttpRequest,
		ctrCMSMostView,
		middlewareAdminAuth.Authenticate,
	))

	//Path public movie
	pathPublicMovie.Get("/list", rtr.handle(
		handler.HttpRequest,
		ctrPublicMovieList,
	))
	pathPublicMovie.Post("/track", rtr.handle(
		handler.HttpRequest,
		ctrPublicTrackMovieViewer,
	))
	pathPublicMovie.Get("/search", rtr.handle(
		handler.HttpRequest,
		ctrPublicMovieSearch,
	))

	//Path user movie
	pathUserMovie.Post("/feedback/:value", rtr.handle(
		handler.HttpRequest,
		ctrUserMovieVote,
		middlewareUserAuth.Authenticate,
	))
}

func NewRouter(cfg *config.Config, fiber fiber.Router) Router {
	return &router{
		cfg:   cfg,
		fiber: fiber,
	}
}
