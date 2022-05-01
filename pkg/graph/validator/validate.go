package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/wascript3r/anomaly/pkg/graph"
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

func (v *Validate) TrapMFCoeffs(c []int, minVal, maxVal *int) error {
	if len(c) > 4 {
		err := v.TrapMFCoeffs(c[:4], minVal, maxVal)
		if err != nil {
			return err
		}
		return v.TrapMFCoeffs(c[4:], minVal, maxVal)
	}

	for i := 0; i < len(c)-1; i++ {
		if c[i] > c[i+1] {
			return graph.InvalidCoeffsError
		}
	}

	if (minVal != nil && c[0] < *minVal) || (maxVal != nil && c[len(c)-1] > *maxVal) {
		return graph.CoeffsOutOfRangeError
	}
	return nil
}
