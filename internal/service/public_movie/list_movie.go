package publicmovie

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
	"strings"
)

func (c *cmsMovieService) List(ctx context.Context, param presentation.ReqPublicMovieList) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServicePublivMovieList"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.public_movie.list", nil)
	defer span.End()

	movies, count, err := c.repoMovie.List(ctx, entity.MetaPagination{
		Page:  param.Page,
		Limit: param.Limit,
	}, []string{"id", "title", "genre_ids", "minutes_duration", "view_number", "created_at", "updated_at"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("fetch list movie got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusInternalServerError)
	}

	var (
		response   []presentation.RespPublicMovieList
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

		response = append(response, presentation.RespPublicMovieList{
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

	return *appctx.NewResponse().WithData(response).WithMeta(appctx.MetaData{
		Page:       param.Page,
		Limit:      param.Limit,
		TotalPage:  helper.PageCalculate(count, param.Limit),
		TotalCount: count,
	})

}
