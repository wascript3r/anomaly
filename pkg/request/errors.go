package request

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

var (
	// Error codes

	InvalidInputError = errcode.InvalidInputError
	UnknownError      = errcode.UnknownError

	CannotParseTimestampError = errcode.New(
		"cannot_parse_timestamp",
		errors.New("cannot parse timestamp"),
	)
)
