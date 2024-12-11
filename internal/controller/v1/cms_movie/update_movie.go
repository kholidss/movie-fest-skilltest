package cmsmovie

import (
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	cmsMovieSvc "github.com/kholidss/movie-fest-skilltest/internal/service/cms_movie"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"

	"github.com/gofiber/fiber/v2"
)

type cmsMovieUpdate struct {
	svcCMSMovie cmsMovieSvc.CMSMovieService
}

func NewCMSMovieUpdate(svcCMSMovie cmsMovieSvc.CMSMovieService) contract.Controller {
	return &cmsMovieUpdate{
		svcCMSMovie: svcCMSMovie,
	}
}

func (cx *cmsMovieUpdate) Serve(xCtx appctx.Data) appctx.Response {
	var (
		authInfo  = helper.GetUserAuthDataFromFiberCtx(xCtx.FiberCtx)
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("CMSV1CreateMovie"),
			logger.Any("X-Request-ID", requestID),
		)
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.cms.update_movie_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the JSON request body
	payload, err := cx.parse(xCtx.FiberCtx)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse form-data payload got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	lf.Append(logger.Any("payload.title", payload.Title))
	lf.Append(logger.Any("payload.genre_ids", payload.GenreIDS))
	lf.Append(logger.Any("payload.minutes_duration", payload.MinutesDuration))
	lf.Append(logger.Any("payload.artists", payload.Artists))
	lf.Append(logger.Any("payload.watch_url", payload.WatchURL))

	// Validate payload
	err = cx.validate(payload)
	if err != nil {
		logger.WarnWithContext(ctx, "payload got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	// Validate file
	errFile := cx.validateFile(payload)
	if len(errFile) > 0 {
		logger.WarnWithContext(ctx, "error file validation", lf...)
		return *appctx.NewResponse().WithError(errFile).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := cx.svcCMSMovie.Update(ctx, authInfo, payload)
	return rsp
}
