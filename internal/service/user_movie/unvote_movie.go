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
	"time"
)

func (u *userMovieService) UnVote(ctx context.Context, authData presentation.UserAuthData, payload presentation.ReqUserMovieVote) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceUserMovieUnVote"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
			logger.Any("payload.movie_id", payload.MovieID),

			logger.Any("actor.user_id", authData.UserID),
			logger.Any("actor.email", authData.Email),
			logger.Any("actor.full_name", authData.FullName),
			logger.Any("actor.user_agent", authData.UserAgent),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.user_movie.unvote", nil)
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

	//Find user already un-vote
	alreadyUnVote, errTrx := u.repoMovieVote.FindOneWithForUpdate(ctx, entity.MovieVote{
		MovieID: payload.MovieID,
		UserID:  authData.UserID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find one movie user vote got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	if alreadyUnVote == nil {
		errTrx = errors.New("already un-vote")
		logger.WarnWithContext(ctx, "user already un-vote this movie", lf...)
		return *appctx.NewResponse().WithCode(http.StatusUnprocessableEntity).WithMessage("User already un-vote this movie")
	}

	lf.Append(logger.Any("movie.id", movie.ID))
	lf.Append(logger.Any("movie.title", movie.Title))
	lf.Append(logger.Any("movie.view_number", movie.ViewNumber))
	lf.Append(logger.Any("movie.vote_number", movie.VoteNumber))

	var (
		resultOfVoteNumber = func() int {
			if movie.VoteNumber > 0 {
				return movie.VoteNumber - 1
			}
			return 0
		}()
		tnow = time.Now()
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

	errTrx = u.repoMovieVote.Update(ctx, entity.MovieVote{
		IsDeleted: true,
		DeletedAt: &tnow,
	}, entity.MovieVote{
		ID: alreadyUnVote.ID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("soft delete user movie vote data got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	//commit db transaction
	_ = tx.Commit()

	//async store history
	go func() {
		_ = u.storeActionHistory(context.Background(), entity.ActionHistory{
			ID:             uuid.New().String(),
			Name:           fmt.Sprintf(consts.ActionHistoryUnVoteMovie, movie.ID, movie.Title),
			IdentifierID:   authData.UserID,
			IdentifierName: consts.RoleEntityUser,
			UserAgent:      authData.UserAgent,
		})
	}()

	logger.InfoWithContext(ctx, "success un-voted movie", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusOK).
		WithMessage("Success un-voted movie").
		WithData(presentation.RespPublicTrackMovieViewer{
			ID:    movie.ID,
			Title: movie.Title,
		})

}
