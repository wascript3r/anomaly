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
	MinVal   *int
	MaxVal   *int
	TrapMFs  []*TrapMF
}

type TrapMF struct {
	ID     int
	Name   string
	Coeffs []int
}

type FullTrapMF struct {
	ID     int
	Name   string
	Coeffs []int
	MinVal *int
	MaxVal *int
}

type RuleText struct {
	ID     int
	Inputs []string
	Output int
}

type Rule struct {
	ID     int
	TFIDs  []int
	Output int
}
