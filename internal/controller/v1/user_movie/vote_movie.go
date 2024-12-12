package usermovie

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	userMovieSvc "github.com/kholidss/movie-fest-skilltest/internal/service/user_movie"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
)

type userVoteMovie struct {
	svcUserMovie userMovieSvc.UserMovieService
}

func NewUserVoteMovie(svcUserMovie userMovieSvc.UserMovieService) contract.Controller {
	return &userVoteMovie{
		svcUserMovie: svcUserMovie,
	}
}

func (ux *userVoteMovie) Serve(xCtx appctx.Data) appctx.Response {
	var (
		authInfo  = helper.GetUserAuthDataFromFiberCtx(xCtx.FiberCtx)
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("UserMovieV1Vote"),
			logger.Any("X-Request-ID", requestID),

			logger.Any("actor.user_id", authInfo.UserID),
			logger.Any("actor.email", authInfo.Email),
			logger.Any("actor.full_name", authInfo.FullName),
			logger.Any("actor.user_agent", authInfo.UserAgent),
		)

		payload presentation.ReqUserMovieVote
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.user_movie.vote_v1", nil)
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

	payload.Value = xCtx.FiberCtx.Params("value")

	lf.Append(logger.Any("payload.movie_id", payload.MovieID))

	// Validate payload
	err = ux.validate(payload)
	if err != nil {
		logger.WarnWithContext(ctx, "payload got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	if payload.Value == "vote" {
		return ux.svcUserMovie.Vote(ctx, authInfo, payload)
	}

	return ux.svcUserMovie.UnVote(ctx, authInfo, payload)
}
