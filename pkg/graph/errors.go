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

	CoeffsOutOfRangeError = errcode.New(
		"coeffs_out_of_range",
		errors.New("coefficients are out of range"),
	)

	InvalidCoeffsDimError = errcode.New(
		"invalid_coeffs_dim",
		errors.New("coefficients must have the same dimension"),
	)

	NotFoundError = errcode.New(
		"not_found",
		errors.New("not found"),
	)

	InvalidOutputIDError = errcode.New(
		"invalid_output_id",
		errors.New("invalid output ID"),
	)
)
