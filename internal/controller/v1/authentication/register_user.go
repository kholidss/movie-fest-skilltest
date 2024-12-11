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
	"github.com/kholidss/movie-fest-skilltest/pkg/masker"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"

	"github.com/gofiber/fiber/v2"
)

type registerUser struct {
	svcAuthentication authentication.AuthenticationService
}

func NewRegisterUser(svcAuthentication authentication.AuthenticationService) contract.Controller {
	return &registerUser{
		svcAuthentication: svcAuthentication,
	}
}

func (r *registerUser) Serve(xCtx appctx.Data) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("AuthV1RegisterUser"),
			logger.Any("X-Request-ID", requestID),
		)
		payload presentation.ReqRegisterUser
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.auth.register_user_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the JSON request body
	err := xCtx.FiberCtx.BodyParser(&payload)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse json payload got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	lf.Append(logger.Any("payload.full_name", payload.FullName))
	lf.Append(logger.Any("payload.email", payload.Email))
	lf.Append(logger.Any("payload.password", masker.Censored(payload.Password, "*")))

	// Validate payload
	err = r.validate(payload)
	if err != nil {
		logger.WarnWithContext(ctx, "payload got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := r.svcAuthentication.RegisterUser(ctx, payload)
	return rsp
}
