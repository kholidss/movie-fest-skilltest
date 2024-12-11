package middleware

import (
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/httpclient"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"github.com/kholidss/movie-fest-skilltest/pkg/util"
	"github.com/spf13/cast"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type adminAuth struct {
	cfg      *config.Config
	repoUser repositories.UserRepository
}

func NewAdminAuthMiddleware(cfg *config.Config, repoUser repositories.UserRepository) *adminAuth {
	return &adminAuth{
		cfg:      cfg,
		repoUser: repoUser,
	}
}

func (a *adminAuth) Authenticate(xCtx *fiber.Ctx) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx)
		lf        = logger.NewFields(
			logger.EventName("AdminAuthMiddleware"),
			logger.Any("X-Request-ID", requestID),
		)
	)

	ctx, span := tracer.NewSpan(xCtx.Context(), "middleware.admin_auth", nil)
	defer span.End()

	headerAuth := xCtx.Get(httpclient.Authorization)

	if headerAuth == "" {
		logger.WarnWithContext(ctx, "authorization header is missing", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	if !strings.HasPrefix(headerAuth, "Bearer ") {
		logger.WarnWithContext(ctx, "Authorization header is missing 'Bearer' prefix", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	token := strings.TrimPrefix(headerAuth, "Bearer ")

	// Validate the JWT token using the public key
	_, claims, err := util.ValidateJWT(token, a.cfg.AppConfig.APPPublicKey)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.WarnWithContext(ctx, "failed to validate user JWT token", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	// Parse user claims
	userClaim, ok := claims["user"].(map[string]any)
	if !ok {
		logger.WarnWithContext(ctx, "user claims not found", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	var (
		userID = cast.ToString(userClaim["user_id"])
	)

	//Fetch user data
	user, err := a.repoUser.FindOne(ctx, entity.User{
		ID: userID,
	}, []string{"id", "full_name", "email", "entity"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("fetch user data got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}
	if user == nil {
		logger.WarnWithContext(ctx, "got null result of user data", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	if user.Entity != consts.RoleEntityAdmin {
		logger.WarnWithContext(ctx, "forbidden access, entity must be admin", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusForbidden).WithMessage(consts.MsgAPIForbidden)
	}

	//Inject user auth data to fiber context
	xCtx.Locals(consts.CtxKeyUserAuthData, presentation.UserAuthData{
		UserID:      user.ID,
		AccessToken: token,
		FullName:    user.FullName,
		Email:       user.Email,
		Entity:      strings.ToLower(user.Entity),
	})

	return *appctx.NewResponse().WithCode(fiber.StatusOK)
}
