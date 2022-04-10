package main

import (
	"fmt"

	"github.com/wascript3r/anomaly/pkg/fuzzy"
	_fuzzyUcase "github.com/wascript3r/anomaly/pkg/fuzzy/usecase"
)

func main() {
	fuzzyUcase, err := _fuzzyUcase.New(
		[]_fuzzyUcase.MembershipFunc{
			_fuzzyUcase.NewTrapMF(6, 9, 11, 14),
			_fuzzyUcase.NewTrapMF(11, 14, 17, 19),
			_fuzzyUcase.NewTrapMF(17, 19, 22, 23),
			_fuzzyUcase.NewCombinedMF(
				_fuzzyUcase.NewTrapMF(21, 23, 24, 24),
				_fuzzyUcase.NewTrapMF(0, 0, 5, 8),
			),
		},
		[]_fuzzyUcase.MembershipFunc{
			_fuzzyUcase.NewTrapMF(1, 1, 4, 6),
			_fuzzyUcase.NewTrapMF(4, 6, 7, 8),
		},
		[]_fuzzyUcase.MembershipFunc{
			_fuzzyUcase.NewTrapMF(0, 0, 5, 8),
			_fuzzyUcase.NewTrapMF(5, 9, 15, 20),
			_fuzzyUcase.NewCombinedMF(
				_fuzzyUcase.NewTrapMF(14, 21, 40, 40),
				func(x float64) float64 {
					if x >= 40 {
						return 1
					}
					return 0
				},
			),
		},
		[]_fuzzyUcase.MembershipFunc{
			_fuzzyUcase.NewTrapMF(0, 0, 30, 50),
			_fuzzyUcase.NewTrapMF(30, 55, 80, 100),
			_fuzzyUcase.NewCombinedMF(
				_fuzzyUcase.NewTrapMF(82, 105, 120, 120),
				func(x float64) float64 {
					if x >= 120 {
						return 1
					}
					return 0
				},
			),
		},
		_fuzzyUcase.NewProbability(
			[]_fuzzyUcase.MembershipFunc{
				_fuzzyUcase.NewTrapMF(0, 0, 20, 38),
				_fuzzyUcase.NewTrapMF(25, 40, 50, 60),
				_fuzzyUcase.NewTrapMF(50, 60, 70, 80),
				_fuzzyUcase.NewTrapMF(72, 85, 100, 100),
			},
			1, 100, 0.1,
		),
		[][]_fuzzyUcase.RuleValue{
			{1, 1, 1, 1, 1},
			{1, 1, 1, 2, 1},
			{1, 1, 1, 3, 2},
			{1, 1, 2, 1, 2},
			{1, 1, 2, 2, 2},
			{1, 1, 2, 3, 3},
			{1, 1, 3, 1, 3},
			{1, 1, 3, 2, 3},
			{1, 1, 3, 3, 4},
			{1, 2, 1, 1, 1},
			{1, 2, 1, 2, 2},
			{1, 2, 1, 3, 3},
			{1, 2, 2, 1, 2},
			{1, 2, 2, 2, 3},
			{1, 2, 2, 3, 4},
			{1, 2, 3, 1, 3},
			{1, 2, 3, 2, 4},
			{1, 2, 3, 3, 4},

			{2, 1, 1, 1, 1},
			{2, 1, 1, 2, 1},
			{2, 1, 1, 3, 2},
			{2, 1, 2, 1, 2},
			{2, 1, 2, 2, 2},
			{2, 1, 2, 3, 2},
			{2, 1, 3, 1, 3},
			{2, 1, 3, 2, 3},
			{2, 1, 3, 3, 4},
			{2, 2, 1, 1, 1},
			{2, 2, 1, 2, 1},
			{2, 2, 1, 3, 3},
			{2, 2, 2, 1, 2},
			{2, 2, 2, 2, 3},
			{2, 2, 2, 3, 4},
			{2, 2, 3, 1, 3},
			{2, 2, 3, 2, 4},
			{2, 2, 3, 3, 4},

			{3, 1, 1, 1, 1},
			{3, 1, 1, 2, 2},
			{3, 1, 1, 3, 2},
			{3, 1, 2, 1, 2},
			{3, 1, 2, 2, 2},
			{3, 1, 2, 3, 3},
			{3, 1, 3, 1, 3},
			{3, 1, 3, 2, 4},
			{3, 1, 3, 3, 4},
			{3, 2, 1, 1, 1},
			{3, 2, 1, 2, 2},
			{3, 2, 1, 3, 2},
			{3, 2, 2, 1, 2},
			{3, 2, 2, 2, 2},
			{3, 2, 2, 3, 3},
			{3, 2, 3, 1, 3},
			{3, 2, 3, 2, 3},
			{3, 2, 3, 3, 4},

			{4, 1, 1, 1, 2},
			{4, 1, 1, 2, 3},
			{4, 1, 1, 3, 3},
			{4, 1, 2, 1, 3},
			{4, 1, 2, 2, 3},
			{4, 1, 2, 3, 4},
			{4, 1, 3, 1, 4},
			{4, 1, 3, 2, 4},
			{4, 1, 3, 3, 4},
			{4, 2, 1, 1, 2},
			{4, 2, 1, 2, 2},
			{4, 2, 1, 3, 3},
			{4, 2, 2, 1, 3},
			{4, 2, 2, 2, 4},
			{4, 2, 2, 3, 4},
			{4, 2, 3, 1, 4},
			{4, 2, 3, 2, 4},
			{4, 2, 3, 3, 4},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(fuzzyUcase.CalcResult(&fuzzy.Model{
		DayTime:   23,
		WeekDay:   7,
		IMSICalls: 45,
		MSCCalls:  15,
	}))
}
