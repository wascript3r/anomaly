package graph

type Validate interface {
	RawRequest(s interface{}) error
	TrapMFCoeffs(c []int) bool
}
