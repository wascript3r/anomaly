package usecase

import (
	"context"
	"errors"
	"time"

	fuzzyLib "github.com/kczimm/fuzzy"
	"github.com/wascript3r/anomaly/pkg/domain"
	"github.com/wascript3r/anomaly/pkg/fuzzy"
	"github.com/wascript3r/anomaly/pkg/graph"
)

const (
	TotalInputCount = 4
)

var (
	ErrInvalidRuleDimensions = errors.New("invalid rule dimensions")
	ErrInvalidRuleValue      = errors.New("invalid rule value")
)

type RuleValue int
type MembershipFunc func(float64) float64

func NewTrapMF(a, b, c, d float64) MembershipFunc {
	return MembershipFunc(fuzzyLib.NewTrapMF(a, b, c, d))
}

func NewCombinedMF(mf1, mf2 MembershipFunc) MembershipFunc {
	return func(x float64) float64 {
		return mf1(x) + mf2(x)
	}
}

func constantMF(y float64) func(float64) float64 {
	return func(_ float64) float64 {
		return y
	}
}

type Probability struct {
	MFs              []MembershipFunc
	MinX, MaxX, Step float64
	x                []float64
	sets             []fuzzyLib.Set
}

func arange(start, end, step float64) []float64 {
	var ret []float64
	for i := start; i <= end; i += step {
		ret = append(ret, i)
	}
	return ret
}

func NewProbability(mfs []MembershipFunc, minX, maxX, step float64) *Probability {
	x := arange(minX, maxX, step)
	sets := make([]fuzzyLib.Set, len(mfs))
	for i, mf := range mfs {
		sets[i] = fuzzyLib.NewFuzzySetFromMF(x, fuzzyLib.MembershipFunc(mf))
	}

	return &Probability{
		MFs:  mfs,
		MinX: minX,
		MaxX: maxX,
		Step: step,
		x:    x,
		sets: sets,
	}
}

type Usecase struct {
	graphRepo  graph.Repository
	ctxTimeout time.Duration

	dayTimeMFs   []MembershipFunc
	weekDayMFs   []MembershipFunc
	imsiCallsMFs []MembershipFunc
	mscCallsMFs  []MembershipFunc
	probability  *Probability
	rules        map[RuleValue][][]RuleValue
}

func (u *Usecase) validateRules(rules [][]RuleValue) error {
	if len(rules) != len(u.dayTimeMFs)*len(u.weekDayMFs)*len(u.imsiCallsMFs)*len(u.mscCallsMFs) {
		return ErrInvalidRuleDimensions
	}

	for _, rule := range rules {
		if len(rule) != TotalInputCount+1 {
			return ErrInvalidRuleDimensions
		}
		mfs := [][]MembershipFunc{u.dayTimeMFs, u.weekDayMFs, u.imsiCallsMFs, u.mscCallsMFs, u.probability.MFs}
		for i, val := range rule {
			if val <= 0 || int(val) > len(mfs[i]) {
				return ErrInvalidRuleValue
			}
		}
	}

	return nil
}

func formatRules(rules [][]RuleValue) map[RuleValue][][]RuleValue {
	ret := make(map[RuleValue][][]RuleValue)
	for _, rule := range rules {
		probVal := rule[len(rule)-1]
		slc, ok := ret[probVal]
		if !ok {
			slc = make([][]RuleValue, 0)
			ret[probVal] = slc
		}
		ret[probVal] = append(slc, rule[:len(rule)-1])
	}
	return ret
}

func New(gr graph.Repository, t time.Duration, rules [][]RuleValue) (*Usecase, error) {
	u := &Usecase{
		graphRepo:  gr,
		ctxTimeout: t,
	}

	err := u.updateConfig(context.Background())
	if err != nil {
		return nil, err
	}

	err = u.validateRules(rules)
	if err != nil {
		return nil, err
	}
	u.rules = formatRules(rules)

	return u, nil
}

func parseTrapMFs(g *domain.Graph) []MembershipFunc {
	var ret []MembershipFunc
	for i, t := range g.TrapMFs {
		var trapMF MembershipFunc
		if (len(t.Coeffs)) == 4 {
			trapMF = NewTrapMF(float64(t.Coeffs[0]), float64(t.Coeffs[1]), float64(t.Coeffs[2]), float64(t.Coeffs[3]))
		} else if (len(t.Coeffs)) == 8 {
			trapMF = NewCombinedMF(
				NewTrapMF(float64(t.Coeffs[0]), float64(t.Coeffs[1]), float64(t.Coeffs[2]), float64(t.Coeffs[3])),
				NewTrapMF(float64(t.Coeffs[4]), float64(t.Coeffs[5]), float64(t.Coeffs[6]), float64(t.Coeffs[7])),
			)
		} else {
			panic("Invalid trapMF")
		}

		if g.Infinite && i == len(g.TrapMFs)-1 {
			trapMF = NewCombinedMF(
				trapMF,
				func(x float64) float64 {
					if x >= float64(t.Coeffs[len(t.Coeffs)-1]) {
						return 1
					}
					return 0
				},
			)
		}

		ret = append(ret, trapMF)
	}

	return ret
}

func (u *Usecase) updateConfig(ctx context.Context) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	gs, err := u.graphRepo.GetAll(c)
	if err != nil {
		return err
	}

	if len(gs) != 5 {
		panic("Invalid graph count")
	}

	u.dayTimeMFs = parseTrapMFs(gs[domain.DayTimeGraphType-1])
	u.weekDayMFs = parseTrapMFs(gs[domain.WeekDayGraphType-1])
	u.imsiCallsMFs = parseTrapMFs(gs[domain.IMSICallsGraphType-1])
	u.mscCallsMFs = parseTrapMFs(gs[domain.MSCCallsGraphType-1])
	u.probability = NewProbability(
		parseTrapMFs(gs[domain.ProbabilityGraphType-1]),
		1, 100, 0.1,
	)

	return nil
}

func calcMFValues(mfs []MembershipFunc, x float64) []float64 {
	vals := make([]float64, len(mfs))
	for i, mf := range mfs {
		vals[i] = mf(x)
	}
	return vals
}

func minSlice(vals []float64) float64 {
	min := vals[0]
	for _, val := range vals {
		if val < min {
			min = val
		}
	}
	return min
}

func maxSlice(vals []float64) float64 {
	max := vals[0]
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max
}

func (u *Usecase) implication(m *fuzzy.Model) map[RuleValue][]float64 {
	dayTimeVals := calcMFValues(u.dayTimeMFs, m.DayTime)
	weekDayVals := calcMFValues(u.weekDayMFs, m.WeekDay)
	imsiCallsVals := calcMFValues(u.imsiCallsMFs, m.IMSICalls)
	mscCallsVals := calcMFValues(u.mscCallsMFs, m.MSCCalls)

	ruleVals := make(map[RuleValue][]float64)
	for probVal, rules := range u.rules {
		ruleVals[probVal] = make([]float64, len(rules))

		for i, rule := range rules {
			ruleVals[probVal][i] = minSlice([]float64{
				dayTimeVals[rule[0]-1],
				weekDayVals[rule[1]-1],
				imsiCallsVals[rule[2]-1],
				mscCallsVals[rule[3]-1],
			})
		}
	}

	return ruleVals
}

func (u *Usecase) aggregation(ruleVals map[RuleValue][]float64) fuzzyLib.Set {
	aggregated := fuzzyLib.NewEmptySet()

	for probVal, vals := range ruleVals {
		probLevel := maxSlice(vals)
		probActivationSet := fuzzyLib.NewFuzzySetFromMF(u.probability.x, constantMF(probLevel)).Intersection(u.probability.sets[probVal-1])
		aggregated = aggregated.Union(probActivationSet)
	}

	return aggregated
}

func (u *Usecase) defuzzification(aggregated fuzzyLib.Set) float64 {
	return aggregated.Centroid()
}

func (u *Usecase) CalcResult(ctx context.Context, m *fuzzy.Model) (float64, error) {
	if err := u.updateConfig(ctx); err != nil {
		return 0, err
	}

	ruleVals := u.implication(m)
	aggregated := u.aggregation(ruleVals)
	return u.defuzzification(aggregated), nil
}
