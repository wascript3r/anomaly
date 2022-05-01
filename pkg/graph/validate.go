package graph

type Validate interface {
	RawRequest(s interface{}) error
	TrapMFCoeffs(c []int, minVal, maxVal *int) error
}
