package cmsmovie

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/pkg/validator"
)

func (cx *cmsMovieCreate) validate(payload presentation.ReqCMSCreateMovie) error {
	rules := []*validation.FieldRules{
		// Title
		validation.Field(&payload.Title, validation.Required, validation.Length(1, 255)),

		// GenreIDS
		validation.Field(&payload.GenreIDS, validation.Required, validation.Each(is.UUID), validation.By(validator.ValidateNoDuplicate)),

		// Description
		validation.Field(&payload.Description, validation.Required),

		// MinutesDuration
		validation.Field(&payload.MinutesDuration, validation.Required),

		// Artists
		validation.Field(&payload.Artists, validation.Required),

		// WatchURL
		validation.Field(&payload.WatchURL, validation.Required, is.URL),

		// FileImage
		validation.Field(&payload.FileImage, validation.Required),
	}

	err := validation.ValidateStruct(&payload, rules...)
	ve, ok := err.(validation.Errors)
	if !ok {
		ve = make(validation.Errors)
	}

	if len(ve) == 0 {
		return nil
	}

	return ve
}
