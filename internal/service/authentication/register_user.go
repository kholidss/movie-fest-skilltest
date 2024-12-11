package authentication

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/pkg/cipher"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/masker"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"net/http"
)

func (a authenticationService) RegisterUser(ctx context.Context, payload presentation.ReqRegisterUser) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceAuthRegisterUser"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.auth.register_user", nil)
	defer span.End()

	lf.Append(logger.Any("payload.email", payload.Email))
	lf.Append(logger.Any("payload.full_name", payload.FullName))
	lf.Append(logger.Any("payload.password", masker.Censored(payload.Password, "*")))

	//Find exist email
	user, err := a.repoUser.FindOne(ctx, entity.User{
		Email: payload.Email,
	}, []string{"id", "full_name", "entity", "created_at", "updated_at"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find exist user by email got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	if user != nil {
		lf.Append(logger.Any("exist_user.id", user.ID))
		lf.Append(logger.Any("exist_user.full_name", user.FullName))
		lf.Append(logger.Any("exist_user.entity", user.Entity))
		lf.Append(logger.Any("exist_user.created_at", user.CreatedAt))
		lf.Append(logger.Any("exist_user.updated_at", user.UpdatedAt))

		logger.WarnWithContext(ctx, "user by email already registered", lf...)
		return *appctx.NewResponse().WithCode(http.StatusUnprocessableEntity).WithMessage("Email already registered")
	}

	encxPassword, err := cipher.EncryptAES256(payload.Password, a.cfg.AppConfig.AppPasswordSecret)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("encrypt user password with aes256 method got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	var (
		userID = uuid.New().String()
	)

	lf.Append(logger.Any("result.user_id", userID))

	//Store user data
	err = a.repoUser.Store(ctx, entity.User{
		ID:       userID,
		Email:    payload.Email,
		FullName: payload.FullName,
		Password: encxPassword,
		Entity:   consts.RoleEntityUser,
	})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store user data got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	logger.InfoWithContext(ctx, "success register user", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusCreated).
		WithMessage("Success register user").
		WithData(presentation.RespRegisterUser{
			UserID:   userID,
			FullName: payload.FullName,
			Email:    payload.Email,
		})

}
