package usermovie

import (
	"context"
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
)

func (u *userMovieService) storeActionHistory(ctx context.Context, history entity.ActionHistory) error {
	var (
		lf = logger.NewFields(
			logger.EventName("SubServiceUserMovieStoreHistory"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "sub_service.user_movie.store_history", nil)
	defer span.End()

	lf.Append(logger.Any("history.id", history.ID))
	lf.Append(logger.Any("history.name", history.Name))
	lf.Append(logger.Any("his.identifier_id", history.IdentifierID))
	lf.Append(logger.Any("history.identifier_name", history.IdentifierName))
	lf.Append(logger.Any("history.user_agent", history.UserAgent))

	err := u.repoActionHistory.Store(ctx, history)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store action history got error: %v", err), lf...)
		return err
	}

	logger.InfoWithContext(ctx, "success store action history data", lf...)
	return nil

}
