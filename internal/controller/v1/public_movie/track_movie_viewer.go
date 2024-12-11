package publicmovie

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

type publicTrackMovieViewer struct {
	publicMovie publicmovie.PublicMovieService
}

func NewPublicTrackMovieViewer(publicMovie publicmovie.PublicMovieService) contract.Controller {
	return &publicTrackMovieViewer{
		publicMovie: publicMovie,
	}
}

func (px *publicTrackMovieViewer) Serve(xCtx appctx.Data) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("PublicMovieV1TrackViewer"),
			logger.Any("X-Request-ID", requestID),
		)

		payload presentation.ReqPublicTrackMovieViewer
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.public_movie.track_viewer_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the JSON request body
	err := xCtx.FiberCtx.BodyParser(&payload)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse json payload got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	lf.Append(logger.Any("param.movie_id", payload.MovieID))

	// Validate payload
	err = px.validate(payload)
	if err != nil {
		logger.WarnWithContext(ctx, "payload got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := px.publicMovie.Track(ctx, payload)
	return rsp
}
