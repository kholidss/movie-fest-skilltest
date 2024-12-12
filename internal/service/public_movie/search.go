package publicmovie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"net/http"
	"strings"
)

func (c *publicMovieService) Search(ctx context.Context, param presentation.ReqPublicMovieSearch) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServicePublicMovieSearch"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.public_movie.search", nil)
	defer span.End()

	lf.Append(logger.Any("param.like[title]", param.LikeTitle))
	lf.Append(logger.Any("param.like[description]", param.LikeDescription))
	lf.Append(logger.Any("param.like[artists]", param.LikeArtist))
	lf.Append(logger.Any("param.like[watch_url]", param.LikeWatchURL))
	lf.Append(logger.Any("param.equal[genre_id]", param.EqualGenreID))
	lf.Append(logger.Any("param.equal[minutes_duration]", param.EqualMinutesDuration))

	where := struct {
		LikeTitle            string `db:"like_title,omitempty"`
		LikeDescription      string `db:"like_description,omitempty"`
		LikeArtist           string `db:"like_artist,omitempty"`
		LikeWatchURL         string `db:"like_watch_url,omitempty"`
		LikeGenreID          string `db:"like_genre_ids,omitempty"`
		EqualMinutesDuration string `db:"minutes_duration,omitempty"`
	}{
		LikeTitle:            param.LikeTitle,
		LikeDescription:      param.LikeDescription,
		LikeArtist:           param.LikeArtist,
		LikeWatchURL:         param.LikeWatchURL,
		LikeGenreID:          param.EqualGenreID,
		EqualMinutesDuration: param.EqualMinutesDuration,
	}

	movies, count, err := c.repoMovie.ListWithLike(ctx, entity.MetaPagination{
		Page:  param.Page,
		Limit: param.Limit,
	}, where, []string{"id", "title", "artist", "genre_ids", "minutes_duration", "view_number", "watch_url", "created_at", "updated_at"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("fetch list movie got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	var (
		response   []presentation.RespPublicMovieSearch
		genreCache = make(map[string]entity.Genre)
	)

	for _, vm := range movies {
		currentGenres := strings.Split(vm.GenreIDS, ";")
		genresData := []entity.Genre{}

		for _, cg := range currentGenres {
			if cachedGenre, exists := genreCache[cg]; exists {
				genresData = append(genresData, cachedGenre)
				continue
			}

			// Fetch genre data if not in cache
			ge, err := c.repoGenre.FindOne(ctx, entity.Genre{
				ID: cg,
			}, []string{"id", "name", "slug", "created_at", "updated_at"})
			if err != nil {
				return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
			}

			// Add to genre cache
			if ge != nil {
				genreCache[cg] = *ge
				genresData = append(genresData, *ge)
			}
		}

		// Fetch image data for the movie
		image, err := c.repoBucket.FindOne(ctx, entity.Bucket{
			IdentifierID:   vm.ID,
			IdentifierName: consts.BucketIdentifierImageMovie,
		}, []string{"url"})
		if err != nil {
			return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
		}

		response = append(response, presentation.RespPublicMovieSearch{
			ID:              vm.ID,
			Title:           vm.Title,
			Genres:          genresData,
			MinutesDuration: vm.MinutesDuration,
			ViewNumber:      vm.ViewNumber,
			Artists: func() []string {
				return strings.Split(vm.Artist, ";")
			}(),
			WatchURL: vm.WatchURL,
			ImageURL: image.URL,
		})
	}

	//Async record genre view of number
	if param.EqualGenreID != "" {
		go func() {
			_ = c.recordViewGenre(context.Background(), entity.Genre{
				ID: param.EqualGenreID,
			})
		}()
	}

	logger.InfoWithContext(ctx, "success search movie", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusOK).
		WithMessage("Success search movie").
		WithData(func() any {
			if len(response) > 0 {
				return response
			}
			return []any{}
		}()).WithMeta(appctx.MetaData{
		Page:       param.Page,
		Limit:      param.Limit,
		TotalPage:  helper.PageCalculate(count, param.Limit),
		TotalCount: count,
	})
}

func (c *publicMovieService) recordViewGenre(ctx context.Context, genre entity.Genre) error {
	var (
		lf = logger.NewFields(
			logger.EventName("SubServicePublicMovieRecordViewGenre"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "sub_service.public_movie.record_view_genre", nil)
	defer span.End()

	lf.Append(logger.Any("genre.id", genre.ID))
	lf.Append(logger.Any("genre.name", genre.Name))

	//start db transaction
	tx, err := c.repoMovie.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("start db transaction got error: %v", err), lf...)
		return err
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

	//FindOne with FOR UPDATE to avoid race condition
	existGenre, errTrx := c.repoGenre.FindOneWithForUpdate(ctx, entity.Genre{
		ID: genre.ID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find exist genre got error: %v", errTrx), lf...)
		return errTrx
	}

	if existGenre == nil {
		errTrx = errors.New("genre not found")
		logger.WarnWithContext(ctx, "genre not found", lf...)
		return errTrx
	}

	var (
		resultViewNumber = existGenre.ViewNumber + 1
	)

	//Update increment view of number
	errTrx = c.repoGenre.Update(ctx, entity.Genre{
		ViewNumber: resultViewNumber,
	}, entity.Genre{
		ID: genre.ID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("update increment genre view number got error: %v", errTrx), lf...)
		return errTrx
	}

	_ = tx.Commit()
	logger.InfoWithContext(ctx, "success record genre view of number", lf...)
	return nil

}
