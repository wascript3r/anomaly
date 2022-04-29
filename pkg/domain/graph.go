package domain

type GraphType uint8

const (
	DayTimeGraphType = iota + 1
	WeekDayGraphType
	IMSICallsGraphType
	MSCCallsGraphType
	ProbabilityGraphType
)

type Graph struct {
	ID       int
	Name     string
	Infinite bool
	TrapMFs  []*TrapMF
}

type TrapMF struct {
	ID     int
	Name   string
	Coeffs []int
}
