package authentication

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/pkg/validator"
)

func (r *registerUser) validate(payload presentation.ReqRegisterUser) error {
	rules := []*validation.FieldRules{
		// NIK
		validation.Field(&payload.Email, validation.Required, is.Email),

		// FullName
		validation.Field(&payload.FullName, validation.Required, validation.Length(2, 100), validator.ValidateHumanName()),

		// Password
		validation.Field(&payload.Password, validation.Required, validation.Length(8, 50)),
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

func (l *loginUser) validate(payload presentation.ReqLoginUser) error {
	rules := []*validation.FieldRules{
		// Email
		validation.Field(&payload.Email, validation.Required, is.Email),

		// Password
		validation.Field(&payload.Password, validation.Required),
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

func (l *loginAdmin) validate(payload presentation.ReqLoginUser) error {
	rules := []*validation.FieldRules{
		// Email
		validation.Field(&payload.Email, validation.Required, is.Email),

		// Password
		validation.Field(&payload.Password, validation.Required),
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
