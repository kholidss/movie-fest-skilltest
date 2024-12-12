package publicmovie

import (
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	publicmovie "github.com/kholidss/movie-fest-skilltest/internal/service/public_movie"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"

	"github.com/gofiber/fiber/v2"
)

type publicMovieSearch struct {
	svcPublicMovie publicmovie.PublicMovieService
}

func NewPublicMovieSearch(svcPublicMovie publicmovie.PublicMovieService) contract.Controller {
	return &publicMovieSearch{
		svcPublicMovie: svcPublicMovie,
	}
}

func (px *publicMovieSearch) Serve(xCtx appctx.Data) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("PublicMovieV1Search"),
			logger.Any("X-Request-ID", requestID),
		)

		param presentation.ReqPublicMovieSearch
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.public_movie.search_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the QParam request body
	err := xCtx.FiberCtx.QueryParser(&param)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse query param got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	//Parsing manually filter query param
	param.LikeTitle = xCtx.FiberCtx.Query("like[title]")
	param.LikeDescription = xCtx.FiberCtx.Query("like[description]")
	param.LikeArtist = xCtx.FiberCtx.Query("like[artists]")
	param.LikeWatchURL = xCtx.FiberCtx.Query("like[watch_url]")
	param.EqualMinutesDuration = xCtx.FiberCtx.Query("equal[minutes_duration]")
	param.EqualGenreID = xCtx.FiberCtx.Query("equal[genre_id]")

	param.Page = helper.PageDefaultValue(param.Page)
	param.Limit = helper.LimitDefaultValue(param.Limit)

	lf.Append(logger.Any("param.like[title]", param.LikeTitle))
	lf.Append(logger.Any("param.like[description]", param.LikeDescription))
	lf.Append(logger.Any("param.like[artists]", param.LikeArtist))
	lf.Append(logger.Any("param.like[watch_url]", param.LikeWatchURL))
	lf.Append(logger.Any("param.equal[genre_id]", param.EqualGenreID))
	lf.Append(logger.Any("param.equal[minutes_duration]", param.EqualMinutesDuration))
	lf.Append(logger.Any("param.page", param.Page))
	lf.Append(logger.Any("param.limit", param.Limit))

	// Validate param
	err = px.validate(param)
	if err != nil {
		logger.WarnWithContext(ctx, "param got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := px.svcPublicMovie.Search(ctx, param)
	return rsp
}
