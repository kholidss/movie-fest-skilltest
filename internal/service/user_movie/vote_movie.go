package usermovie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"net/http"
)

func (u *userMovieService) Vote(ctx context.Context, authData presentation.UserAuthData, payload presentation.ReqUserMovieVote) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceUserMovieVote"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
			logger.Any("payload.movie_id", payload.MovieID),

			logger.Any("actor.user_id", authData.UserID),
			logger.Any("actor.email", authData.Email),
			logger.Any("actor.full_name", authData.FullName),
			logger.Any("actor.user_agent", authData.UserAgent),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.user_movie.vote", nil)
	defer span.End()

	//start db transaction
	tx, err := u.repoMovie.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("start db transaction got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}
	txOpt := repositories.WithTransaction(tx)

	var (
		errTrx error
	)

	//always rollback db transaction if got error on store process
	defer func() {
		if errTrx != nil && tx != nil {
			_ = tx.Rollback()
		}
	}()

	//Fetch vote number movie data with FOR UPDATE to avoid race condition
	movie, errTrx := u.repoMovie.FindOneWithForUpdate(ctx, entity.Movie{
		ID: payload.MovieID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find one movie with for update got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	if movie == nil {
		errTrx = errors.New("movie not found")
		logger.WarnWithContext(ctx, "got not found movie data", lf...)
		return *appctx.NewResponse().WithCode(http.StatusNotFound).WithMessage("Movie data not found")
	}

	//Find user already vote
	alreadyVote, errTrx := u.repoMovieVote.FindOneWithForUpdate(ctx, entity.MovieVote{
		MovieID: payload.MovieID,
		UserID:  authData.UserID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find one movie user vote got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	if alreadyVote != nil {
		errTrx = errors.New("already vote")
		logger.WarnWithContext(ctx, "user already vote this movie", lf...)
		return *appctx.NewResponse().WithCode(http.StatusUnprocessableEntity).WithMessage("User already vote this movie")
	}

	lf.Append(logger.Any("movie.id", movie.ID))
	lf.Append(logger.Any("movie.title", movie.Title))
	lf.Append(logger.Any("movie.view_number", movie.ViewNumber))
	lf.Append(logger.Any("movie.vote_number", movie.VoteNumber))

	var (
		resultOfVoteNumber = movie.VoteNumber + 1
	)

	lf.Append(logger.Any("result.vote_number", resultOfVoteNumber))

	errTrx = u.repoMovie.Update(ctx, entity.Movie{
		VoteNumber: resultOfVoteNumber,
	}, entity.Movie{
		ID: payload.MovieID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("update movie view number got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	errTrx = u.repoMovieVote.Store(ctx, entity.MovieVote{
		ID:      uuid.New().String(),
		UserID:  authData.UserID,
		MovieID: movie.ID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store user movie vote data got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	//commit db transaction
	_ = tx.Commit()

	//async store history
	go func() {
		_ = u.storeActionHistory(context.Background(), entity.ActionHistory{
			ID:             uuid.New().String(),
			Name:           fmt.Sprintf(consts.ActionHistoryVoteMovie, movie.ID, movie.Title),
			IdentifierID:   authData.UserID,
			IdentifierName: consts.RoleEntityUser,
			UserAgent:      authData.UserAgent,
		})
	}()

	logger.InfoWithContext(ctx, "success voted movie", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusOK).
		WithMessage("Success voted movie").
		WithData(presentation.RespPublicTrackMovieViewer{
			ID:    movie.ID,
			Title: movie.Title,
		})

}
