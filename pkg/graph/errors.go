package graph

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

var (
	// Error codes

	InvalidInputError = errcode.InvalidInputError
	UnknownError      = errcode.UnknownError

	InvalidCoeffsError = errcode.New(
		"invalid_coeffs",
		errors.New("coefficients must be in increasing order"),
	)

	InvalidCoeffsDimError = errcode.New(
		"invalid_coeffs_dim",
		errors.New("coefficients must have the same dimension"),
	)
)
