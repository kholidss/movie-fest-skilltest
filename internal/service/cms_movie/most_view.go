package cmsmovie

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func (c *cmsMovieService) MostView(ctx context.Context, authData presentation.UserAuthData, param presentation.ReqCMSMostView) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceCMSMovieMostView"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
			logger.Any("actor.id", authData.UserID),
			logger.Any("actor.full_name", authData.FullName),
			logger.Any("actor.email", authData.Email),
			logger.Any("actor.entity", authData.Entity),

			logger.Any("param.value", param.Value),
			logger.Any("param.page", param.Page),
			logger.Any("param.limit", param.Limit),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.cms_movie.most_view", nil)
	defer span.End()

	if param.Value == consts.ValueMovie {
		resp, err := c.mostViewMovie(ctx, param)
		if err != nil {
			tracer.AddSpanError(span, err)
			logger.ErrorWithContext(ctx, fmt.Sprintf("fetch list most view movie got error: %v", err), lf...)
			return resp
		}
		logger.InfoWithContext(ctx, "success get list most view movie", lf...)
		return resp
	}

	resp, err := c.mostViewGenre(ctx, param)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("fetch list most view genre got error: %v", err), lf...)
		return resp
	}
	return resp

}

func (c *cmsMovieService) mostViewMovie(ctx context.Context, param presentation.ReqCMSMostView) (appctx.Response, error) {
	movies, count, err := c.repoMovie.ListMostView(ctx, entity.MetaPagination{
		Page:  param.Page,
		Limit: param.Limit,
	}, []string{"id", "title", "genre_ids", "minutes_duration", "view_number", "created_at", "updated_at"})
	if err != nil {
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError), errors.Wrap(err, "fetch list of movies err")
	}

	var (
		response   presentation.RespCMSMostView
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
				return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError), errors.Wrap(err, "find one genre err")
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
			return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError), errors.Wrap(err, "find one image err")
		}

		response.Movies = append(response.Movies, presentation.RespCMSMostViewMoview{
			ID:              vm.ID,
			Title:           vm.Title,
			Genres:          genresData,
			MinutesDuration: vm.MinutesDuration,
			ViewNumber:      vm.ViewNumber,
			Artist:          vm.Artist,
			WatchURL:        vm.WatchURL,
			ImageURL:        image.URL,
		})
	}

	if response.Movies == nil {
		response.Movies = []presentation.RespCMSMostViewMoview{}
	}

	return *appctx.NewResponse().
		WithCode(http.StatusOK).
		WithMessage("Success get most view movie").
		WithData(response).WithMeta(appctx.MetaData{
		Page:       param.Page,
		Limit:      param.Limit,
		TotalPage:  helper.PageCalculate(count, param.Limit),
		TotalCount: count,
	}), nil
}

func (c *cmsMovieService) mostViewGenre(ctx context.Context, param presentation.ReqCMSMostView) (appctx.Response, error) {
	genres, count, err := c.repoGenre.ListMostView(ctx, entity.MetaPagination{
		Page:  param.Page,
		Limit: param.Limit,
	}, []string{"id", "name", "slug", "view_number", "created_at", "updated_at"})
	if err != nil {
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError), errors.Wrap(err, "fetch list of movies err")
	}

	return *appctx.NewResponse().
		WithCode(http.StatusOK).
		WithMessage("Success get most view genre").
		WithData(presentation.RespCMSMostView{
			Genres: genres,
		}).WithMeta(appctx.MetaData{
		Page:       param.Page,
		Limit:      param.Limit,
		TotalPage:  helper.PageCalculate(count, param.Limit),
		TotalCount: count,
	}), nil
}
