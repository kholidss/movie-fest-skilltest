package authentication

import (
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/internal/service/authentication"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"

	"github.com/gofiber/fiber/v2"
)

type loginUser struct {
	svcAuthentication authentication.AuthenticationService
}

func NewLoginUser(svcAuthentication authentication.AuthenticationService) contract.Controller {
	return &loginUser{
		svcAuthentication: svcAuthentication,
	}
}

func (l *loginUser) Serve(xCtx appctx.Data) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("AuthV1LoginUser"),
			logger.Any("X-Request-ID", requestID),
		)
		payload presentation.ReqLoginUser
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.auth.login_user_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the form data request body
	err := xCtx.FiberCtx.BodyParser(&payload)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse json payload got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	lf.Append(logger.Any("payload.email", payload.Email))

	// Validate payload
	err = l.validate(payload)
	if err != nil {
		logger.WarnWithContext(ctx, "payload got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := l.svcAuthentication.LoginUser(ctx, payload)
	return rsp
}
