package publicmovie

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
)

func (px *publicTrackMovieViewer) validate(payload presentation.ReqPublicTrackMovieViewer) error {
	rules := []*validation.FieldRules{
		// Value
		validation.Field(&payload.MovieID, validation.Required, is.UUID),
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
