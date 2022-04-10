package main

import (
	"fmt"

	"github.com/wascript3r/anomaly/pkg/fuzzy"
	_fuzzyUcase "github.com/wascript3r/anomaly/pkg/fuzzy/usecase"
)

func main() {
	fuzzyUcase, err := _fuzzyUcase.New(
		[]_fuzzyUcase.MembershipFunc{
			_fuzzyUcase.NewTrapMF(28, 28, 40, 45),
			_fuzzyUcase.NewTrapMF(42, 45, 58, 63),
			_fuzzyUcase.NewTrapMF(59, 63, 77, 77),
		},
		[]_fuzzyUcase.MembershipFunc{
			_fuzzyUcase.NewTrapMF(60, 60, 90, 100),
			_fuzzyUcase.NewTrapMF(90, 100, 140, 155),
			_fuzzyUcase.NewTrapMF(140, 155, 202, 202),
		},
		[]_fuzzyUcase.MembershipFunc{
			_fuzzyUcase.NewTrapMF(80, 80, 100, 110),
			_fuzzyUcase.NewTrapMF(100, 115, 140, 155),
			_fuzzyUcase.NewTrapMF(135, 160, 200, 200),
		},
		_fuzzyUcase.NewProbability(
			[]_fuzzyUcase.MembershipFunc{
				_fuzzyUcase.NewTrapMF(0, 0, 20, 35),
				_fuzzyUcase.NewTrapMF(20, 40, 60, 70),
				_fuzzyUcase.NewTrapMF(60, 70, 100, 100),
			},
			1, 100, 0.1,
		),
		[][]_fuzzyUcase.RuleValue{
			{1, 1, 1, 1},
			{1, 1, 2, 2},
			{1, 1, 3, 3},
			{1, 2, 1, 2},
			{1, 2, 2, 2},
			{1, 2, 3, 2},
			{1, 3, 1, 2},
			{1, 3, 2, 1},
			{1, 3, 3, 1},
			{2, 1, 1, 2},
			{2, 1, 2, 3},
			{2, 1, 3, 3},
			{2, 2, 1, 2},
			{2, 2, 2, 2},
			{2, 2, 3, 3},
			{2, 3, 1, 1},
			{2, 3, 2, 2},
			{2, 3, 3, 2},
			{3, 1, 1, 2},
			{3, 1, 2, 3},
			{3, 1, 3, 3},
			{3, 2, 1, 3},
			{3, 2, 2, 3},
			{3, 2, 3, 3},
			{3, 3, 1, 2},
			{3, 3, 2, 2},
			{3, 3, 3, 2},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(fuzzyUcase.CalcResult(&fuzzy.Model{
		Age:       74,
		MaxHR:     90,
		RestingBP: 140,
	}))
}
