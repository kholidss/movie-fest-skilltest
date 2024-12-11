package authentication

import (
	"context"
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/dto"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/pkg/cipher"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"github.com/kholidss/movie-fest-skilltest/pkg/util"
	"net/http"
	"time"
)

func (a authenticationService) LoginUser(ctx context.Context, payload presentation.ReqLoginUser) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceAuthLoginUser"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.auth.login_user", nil)
	defer span.End()

	lf.Append(logger.Any("payload.email", payload.Email))

	//Find exist Email
	user, err := a.repoUser.FindOne(ctx, entity.User{
		Email: payload.Email,
	}, []string{"id", "email", "full_name", "entity", "password", "created_at", "updated_at"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find exist user by email got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	if user == nil {
		logger.WarnWithContext(ctx, "user email not found", lf...)
		return *appctx.NewResponse().WithCode(http.StatusUnprocessableEntity).WithMessage("Invalid email or password")
	}

	lf.Append(logger.Any("exist_user.id", user.ID))
	lf.Append(logger.Any("exist_user.full_name", user.FullName))
	lf.Append(logger.Any("exist_user.created_at", user.CreatedAt))
	lf.Append(logger.Any("exist_user.updated_at", user.UpdatedAt))

	plainPassword, err := cipher.DecryptAES256(user.Password, a.cfg.AppConfig.AppPasswordSecret)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("decrypt aes256 user password got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	if plainPassword != payload.Password {
		logger.WarnWithContext(ctx, "password doesn't match", lf...)
		return *appctx.NewResponse().WithCode(http.StatusUnprocessableEntity).WithMessage("Invalid email or password")
	}

	token, err := util.GenerateJWT(a.cfg.AppConfig.APPPrivateKey, dto.BuildAuthJWTClaims(
		time.Now().Add(time.Hour*time.Duration(a.cfg.AppConfig.APPTokenUserExpiredInHour)).Unix(),
		user,
	))
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("generate new jwt bearer token user got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	logger.InfoWithContext(ctx, "success login user", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusCreated).
		WithMessage("Success login user").
		WithData(presentation.RespLoginUser{
			UserID:      user.ID,
			AccessToken: token,
			FullName:    user.FullName,
			Email:       user.Email,
		})

}
