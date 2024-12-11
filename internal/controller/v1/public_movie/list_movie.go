package cmsmovie

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	publicmovie "github.com/kholidss/movie-fest-skilltest/internal/service/public_movie"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
)

type publicListMovie struct {
	publicMovie publicmovie.PublicMovieService
}

func NewPublicListMovie(publicMovie publicmovie.PublicMovieService) contract.Controller {
	return &publicListMovie{
		publicMovie: publicMovie,
	}
}

func (px *publicListMovie) Serve(xCtx appctx.Data) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("PublicMovieV1List"),
			logger.Any("X-Request-ID", requestID),
		)

		param presentation.ReqPublicMovieList
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.public_movie.list_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the JSON request body
	err := xCtx.FiberCtx.QueryParser(&param)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse query param got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	param.Page = helper.PageDefaultValue(param.Page)
	param.Limit = helper.LimitDefaultValue(param.Limit)

	lf.Append(logger.Any("param.page", param.Page))
	lf.Append(logger.Any("param.limit", param.Limit))

	rsp := px.publicMovie.List(ctx, param)
	return rsp
}
