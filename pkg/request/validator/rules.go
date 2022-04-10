package validator

import (
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

type rules struct {
	dateTimeFormat string
}

func newRules(dateTimeFormat string) rules {
	return rules{dateTimeFormat}
}

func (r rules) attachTo(goV *validator.Validate) {
	aliases := map[string]string{
		"r_imsi": "lte=20",
		"r_msc":  "lte=20",
	}
	validators := map[string]validator.Func{
		"datetime": r.validateDateTime,
	}

	for k, v := range aliases {
		goV.RegisterAlias(k, v)
	}
	for k, v := range validators {
		goV.RegisterValidation(k, v)
	}
}

func (r rules) validateDateTime(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() != reflect.String {
		return false
	}
	date := field.String()

	_, err := time.Parse(r.dateTimeFormat, date)
	return err == nil
}
