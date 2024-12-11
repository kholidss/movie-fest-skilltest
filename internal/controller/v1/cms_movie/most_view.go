package cmsmovie

import (
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/controller/contract"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	cmsMovieSvc "github.com/kholidss/movie-fest-skilltest/internal/service/cms_movie"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type cmsMostView struct {
	svcCMSMovie cmsMovieSvc.CMSMovieService
}

func NewCMSMostView(svcCMSMovie cmsMovieSvc.CMSMovieService) contract.Controller {
	return &cmsMostView{
		svcCMSMovie: svcCMSMovie,
	}
}

func (cx *cmsMostView) Serve(xCtx appctx.Data) appctx.Response {
	var (
		authInfo  = helper.GetUserAuthDataFromFiberCtx(xCtx.FiberCtx)
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("CMSV1MostView"),
			logger.Any("X-Request-ID", requestID),
		)

		param presentation.ReqCMSMostView
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.cms.most_view_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the JSON request body
	err := xCtx.FiberCtx.QueryParser(&param)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse query param got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	param.Value = strings.ToLower(param.Value)
	param.Page = helper.PageDefaultValue(param.Page)
	param.Limit = helper.LimitDefaultValue(param.Limit)

	lf.Append(logger.Any("param.value", param.Value))
	lf.Append(logger.Any("param.page", param.Page))
	lf.Append(logger.Any("param.limit", param.Limit))

	// Validate param
	err = cx.validate(param)
	if err != nil {
		logger.WarnWithContext(ctx, "param got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := cx.svcCMSMovie.MostView(ctx, authInfo, param)
	return rsp
}
