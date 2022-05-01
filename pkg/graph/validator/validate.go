package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validate struct {
	govalidate *validator.Validate
}

func New() *Validate {
	return &Validate{validator.New()}
}

func (v *Validate) RawRequest(s interface{}) error {
	return v.govalidate.Struct(s)
}

func (v *Validate) TrapMFCoeffs(c []int) bool {
	for i := 0; i < len(c)-1; i++ {
		if c[i] > c[i+1] {
			return false
		}
	}
	return true
}
