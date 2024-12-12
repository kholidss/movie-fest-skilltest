package cmsmovie

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strings"
	"time"
)

func (c *cmsMovieService) Update(ctx context.Context, authData presentation.UserAuthData, payload presentation.ReqCMSUpdateMovie) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceCMSMovieCreate"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
			logger.Any("actor.id", authData.UserID),
			logger.Any("actor.full_name", authData.FullName),
			logger.Any("actor.email", authData.Email),
			logger.Any("actor.entity", authData.Entity),

			logger.Any("payload.title", payload.Title),
			logger.Any("payload.genre_ids", payload.GenreIDS),
			logger.Any("payload.minutes_duration", payload.MinutesDuration),
			logger.Any("payload.artist", payload.Artists),
			logger.Any("payload.watch_url", payload.WatchURL),
		)

		genres           []entity.Genre
		genreIDSNotExist []string
	)

	ctx, span := tracer.NewSpan(ctx, "service.cms_movie.create", nil)
	defer span.End()

	currentMovie, err := c.repoMovie.FindOne(ctx, entity.Movie{
		ID: payload.MovieID,
	}, []string{"id", "title"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find one current movie got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	if currentMovie == nil {
		logger.WarnWithContext(ctx, "got not found movie data", lf...)
		return *appctx.NewResponse().WithCode(http.StatusNotFound).WithMessage("Movie data not found")
	}

	//find exist genres
	for i, v := range payload.GenreIDS {
		ge, err := c.repoGenre.FindOne(ctx, entity.Genre{
			ID: v,
		}, []string{"id", "name", "created_at", "updated_at"})
		if err != nil {
			tracer.AddSpanError(span, err)
			logger.ErrorWithContext(ctx, fmt.Sprintf("find genre_id: %s got error: %v", v, err), lf...)
			return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
		}

		if ge == nil {
			lf.Append(logger.Any(fmt.Sprintf("genre[%v].id not exist", i), v))
			genreIDSNotExist = append(genreIDSNotExist, v)
			continue
		}

		lf.Append(logger.Any(fmt.Sprintf("genre[%v].id", i), v))
		lf.Append(logger.Any(fmt.Sprintf("genre[%v].name", i), ge.Name))
		lf.Append(logger.Any(fmt.Sprintf("genre[%v].slug", i), ge.Slug))

		genres = append(genres, *ge)
	}

	//if any genre ids not exist
	if len(genreIDSNotExist) > 0 {
		logger.WarnWithContext(ctx, "genre_ids not exist", lf...)
		return *appctx.NewResponse().
			WithCode(fiber.StatusUnprocessableEntity).
			WithMessage("Genre IDS not exist").
			WithError([]appctx.ErrorResp{
				{
					Key:      "genre_ids",
					Messages: genreIDSNotExist,
				},
			})
	}

	//action before update
	err = c.beforeUpdate(ctx, payload)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("do before update movie got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	var (
		objectImage *uploader.UploadResult

		fileNameImage, pathImage string
	)

	//upload movie image to CDN Storage
	if payload.FileImage != nil {
		fileNameImage, pathImage = helper.GeneratePathAndFilenameStorage("movie", strings.Split(payload.FileImage.Mimetype, "/")[1])
		obj, err := c.cdn.Put(ctx, pathImage, payload.FileImage.File)
		if err != nil {
			tracer.AddSpanError(span, err)
			logger.ErrorWithContext(ctx, fmt.Sprintf("upload movie image to cdn storage got error: %v", err), lf...)
			return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
		}
		vx, ok := obj.(*uploader.UploadResult)
		if ok {
			objectImage = vx
		}
	}

	jsonStorageBuilder := func(obj *uploader.UploadResult) json.RawMessage {
		if obj == nil {
			return nil
		}
		v, err := json.Marshal(obj)
		if err != nil {
			return nil
		}
		return v
	}

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

	//update movie data
	errTrx = c.repoMovie.Update(ctx, entity.Movie{
		Title:           payload.Title,
		GenreIDS:        strings.Join(payload.GenreIDS, ";"),
		Description:     payload.Description,
		MinutesDuration: payload.MinutesDuration,
		Artist:          strings.Join(payload.Artists, ";"),
		WatchURL:        payload.WatchURL,
	}, entity.Movie{
		ID: payload.MovieID,
	})
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("update movie data got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	gr, _ := errgroup.WithContext(ctx)

	//async store movie genre
	gr.Go(func() error {
		for _, v := range genres {
			err := c.repoMovieGenre.Store(ctx, entity.MovieGenre{
				ID:      uuid.New().String(),
				MovieID: payload.MovieID,
				GenreID: v.ID,
			}, txOpt)
			if err != nil {
				return err
			}
		}
		return nil
	})

	//async store action history
	gr.Go(func() error {
		return c.repoActionHistory.Store(ctx, entity.ActionHistory{
			ID:             uuid.New().String(),
			Name:           fmt.Sprintf(consts.ActionHistoryUpdateMovie, payload.MovieID, payload.Title),
			IdentifierID:   authData.UserID,
			IdentifierName: consts.RoleEntityAdmin,
			UserAgent:      authData.UserAgent,
		}, txOpt)
	})

	//async store bucket file
	if payload.FileImage != nil {
		gr.Go(func() error {
			return c.repoBucket.Store(ctx, entity.Bucket{
				ID:             uuid.New().String(),
				Filename:       fileNameImage,
				IdentifierID:   payload.MovieID,
				IdentifierName: consts.BucketIdentifierImageMovie,
				Mimetype:       payload.FileImage.Mimetype,
				Provider:       strings.ToLower(c.cfg.CDNConfig.Provider),
				URL:            objectImage.URL,
				Path:           pathImage,
				SupportData:    jsonStorageBuilder(objectImage),
			}, txOpt)
		})
	}

	errTrx = gr.Wait()
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store movie genre or action history data got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	//commit the db transaction
	if tx != nil {
		_ = tx.Commit()

	}

	logger.InfoWithContext(ctx, "success admin cms create movie", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusCreated).
		WithMessage("Success update movie").
		WithData(presentation.RespCMSCreateMovie{
			ID:              payload.MovieID,
			Title:           payload.Title,
			Genres:          genres,
			MinutesDuration: payload.MinutesDuration,
			Artists:         payload.Artists,
			WatchURL:        payload.WatchURL,
			ImageURL: func() string {
				if objectImage != nil {
					return objectImage.URL
				}
				return ""
			}(),
			CreatedBy: presentation.CreatedBy{
				ID:       authData.UserID,
				Email:    authData.Email,
				FullName: authData.FullName,
				Entity:   authData.Entity,
			},
		})
}

func (c *cmsMovieService) beforeUpdate(ctx context.Context, payload presentation.ReqCMSUpdateMovie) error {
	//start db transaction
	tx, err := c.repoMovie.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return errors.Wrap(err, "start db transaction err")
	}

	txOpt := repositories.WithTransaction(tx)

	var (
		errTrx error
		tnow   = time.Now()
	)

	defer func() {
		if errTrx != nil && tx != nil {
			_ = tx.Rollback()
		}
	}()

	gr, _ := errgroup.WithContext(ctx)

	//async delete image movie data
	if payload.FileImage != nil {
		gr.Go(func() error {
			return c.repoBucket.Update(ctx, entity.Bucket{
				IsDeleted: true,
				DeletedAt: &tnow,
			}, entity.Bucket{
				IdentifierID:   payload.MovieID,
				IdentifierName: consts.BucketIdentifierImageMovie,
			}, txOpt)
		})
	}

	//async delete movie genre data
	gr.Go(func() error {
		return c.repoMovieGenre.Update(ctx, entity.MovieGenre{
			IsDeleted: true,
			DeletedAt: &tnow,
		}, entity.MovieGenre{
			MovieID: payload.MovieID,
		}, txOpt)
	})

	errTrx = gr.Wait()
	if errTrx != nil {
		return errors.Wrap(errTrx, "delete movie genre or bucket err")
	}

	//commit db transaction
	if tx != nil {
		_ = tx.Commit()
	}

	return nil

}
