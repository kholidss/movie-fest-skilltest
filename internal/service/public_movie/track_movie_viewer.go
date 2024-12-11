package publicmovie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"net/http"
)

func (c *cmsMovieService) Track(ctx context.Context, payload presentation.ReqPublicTrackMovieViewer) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServicePublicTrackMovieViewer"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
			logger.Any("payload.movie_id", payload.MovieID),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.public_movie.track_viewer", nil)
	defer span.End()

	//start db transaction
	tx, err := c.repoMovie.BeginTx(ctx, &sql.TxOptions{
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

	//Fetch movie data with FOR UPDATE to avoid race condition
	movie, errTrx := c.repoMovie.FindOneWithForUpdate(ctx, entity.Movie{
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

	lf.Append(logger.Any("movie.id", movie.ID))
	lf.Append(logger.Any("movie.title", movie.Title))
	lf.Append(logger.Any("movie.view_number", movie.ViewNumber))

	var (
		resultOfViewer = movie.ViewNumber + 1
	)

	lf.Append(logger.Any("result.view_number", resultOfViewer))

	errTrx = c.repoMovie.Update(ctx, entity.Movie{
		ViewNumber: resultOfViewer,
	}, entity.Movie{
		ID: payload.MovieID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("update movie view number got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	//commit db transaction
	_ = tx.Commit()

	logger.InfoWithContext(ctx, "success record movie viewer", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusOK).
		WithMessage("Success record view of movie").
		WithData(presentation.RespPublicTrackMovieViewer{
			ID:    movie.ID,
			Title: movie.Title,
		})

}
