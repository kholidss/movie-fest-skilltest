package cmsmovie

import (
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/pkg/util"
	"strings"
)

func (cx *cmsMovieCreate) validateFile(payload presentation.ReqCMSCreateMovie) []appctx.ErrorResp {
	var (
		errFormat []appctx.ErrorResp

		msgImage []string
	)

	if !util.InArray(payload.FileImage.Mimetype, consts.AllowedMimeTypesKTP) {
		msgImage = append(
			msgImage,
			fmt.Sprintf("invalid mimetype, only allow: %s", strings.Join(consts.AllowedMimeTypesKTP, ",")),
		)
	}

	if payload.FileImage.Size > consts.MaxSizeKTP {
		msgImage = append(msgImage,
			fmt.Sprintf(
				"file size %s is too large, max file size is: %s",
				util.HumanFileSize(float64(payload.FileImage.Size)),
				util.HumanFileSize(float64(consts.MaxSizeKTP))),
		)
	}

	if len(msgImage) > 0 {
		errFormat = append(errFormat, appctx.ErrorResp{
			Key:      "image",
			Messages: msgImage,
		})
	}

	return errFormat
}

func (cx *cmsMovieUpdate) validateFile(payload presentation.ReqCMSUpdateMovie) []appctx.ErrorResp {
	var (
		errFormat []appctx.ErrorResp
		msgImage  []string
	)

	if payload.FileImage == nil {
		return nil
	}

	if !util.InArray(payload.FileImage.Mimetype, consts.AllowedMimeTypesKTP) {
		msgImage = append(
			msgImage,
			fmt.Sprintf("invalid mimetype, only allow: %s", strings.Join(consts.AllowedMimeTypesKTP, ",")),
		)
	}

	if payload.FileImage.Size > consts.MaxSizeKTP {
		msgImage = append(msgImage,
			fmt.Sprintf(
				"file size %s is too large, max file size is: %s",
				util.HumanFileSize(float64(payload.FileImage.Size)),
				util.HumanFileSize(float64(consts.MaxSizeKTP))),
		)
	}

	if len(msgImage) > 0 {
		errFormat = append(errFormat, appctx.ErrorResp{
			Key:      "image",
			Messages: msgImage,
		})
	}

	return errFormat
}
