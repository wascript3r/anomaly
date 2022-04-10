package main

import (
	"fmt"
	"math"

	fuzzyLib "github.com/kczimm/fuzzy"
)

func ConstantMF(y float64) func(float64) float64 {
	return func(_ float64) float64 {
		return y
	}
}

func main() {
	var (
		inputAge       float64 = 61
		inputMaxHr     float64 = 145
		inputRestingBP float64 = 145
	)

	ageLo := fuzzyLib.NewTrapMF(28, 28, 40, 45)
	ageMd := fuzzyLib.NewTrapMF(42, 45, 58, 63)
	ageHi := fuzzyLib.NewTrapMF(59, 63, 77, 77)

	maxHrLo := fuzzyLib.NewTrapMF(60, 60, 90, 100)
	maxHrMd := fuzzyLib.NewTrapMF(90, 100, 140, 155)
	maxHrHi := fuzzyLib.NewTrapMF(140, 155, 202, 202)

	restingBPLo := fuzzyLib.NewTrapMF(80, 80, 100, 110)
	restingBPMd := fuzzyLib.NewTrapMF(100, 115, 140, 155)
	restingBPHi := fuzzyLib.NewTrapMF(135, 160, 200, 200)

	probLo := fuzzyLib.NewTrapMF(0, 0, 20, 35)
	probMd := fuzzyLib.NewTrapMF(20, 40, 60, 70)
	probHi := fuzzyLib.NewTrapMF(60, 70, 100, 100)
	xProb := arange(1, 100, 0.1)
	probLoSet := fuzzyLib.NewFuzzySetFromMF(xProb, probLo)
	probMdSet := fuzzyLib.NewFuzzySetFromMF(xProb, probMd)
	probHiSet := fuzzyLib.NewFuzzySetFromMF(xProb, probHi)

	ageLevelLo := ageLo(inputAge)
	ageLevelMd := ageMd(inputAge)
	ageLevelHi := ageHi(inputAge)

	maxHrLevelLo := maxHrLo(inputMaxHr)
	maxHrLevelMd := maxHrMd(inputMaxHr)
	maxHrLevelHi := maxHrHi(inputMaxHr)

	restingBPLevelLo := restingBPLo(inputRestingBP)
	restingBPLevelMd := restingBPMd(inputRestingBP)
	restingBPLevelHi := restingBPHi(inputRestingBP)

	fmt.Println("Amžius:", ageLevelLo, ageLevelMd, ageLevelHi)
	fmt.Println("Maksimalus širdies dažnis:", maxHrLevelLo, maxHrLevelMd, maxHrLevelHi)
	fmt.Println("Kraujospūdis:", restingBPLevelLo, restingBPLevelMd, restingBPLevelHi)

	// Implication
	restingBPNotLo := math.Max(restingBPLevelMd, restingBPLevelHi)
	ageNotMd := math.Max(ageLevelLo, ageLevelHi)
	maxHrNotLo := math.Max(maxHrLevelMd, maxHrLevelHi)
	ageNotLo := math.Max(ageLevelMd, ageLevelHi)
	restingBPNotHi := math.Max(restingBPLevelLo, restingBPLevelMd)

	rule1 := math.Min(ageLevelLo, math.Min(maxHrLevelLo, restingBPLevelLo))
	rule2 := math.Min(ageLevelLo, math.Min(maxHrLevelHi, restingBPNotLo))
	rule3 := math.Min(ageLevelMd, math.Min(maxHrLevelHi, restingBPLevelLo))
	probLevelLo := math.Max(rule1, math.Max(rule2, rule3))
	probActivationSetLo := fuzzyLib.NewFuzzySetFromMF(xProb, ConstantMF(probLevelLo)).Intersection(probLoSet)

	rule4 := math.Min(ageLevelLo, math.Min(maxHrLevelLo, restingBPLevelMd))
	rule5 := math.Min(ageNotMd, maxHrNotLo)
	rule6 := math.Min(ageLevelLo, math.Min(maxHrLevelHi, restingBPLevelLo))
	rule7 := math.Min(ageNotLo, math.Min(maxHrLevelLo, restingBPLevelLo))
	rule8 := math.Min(ageLevelMd, math.Min(maxHrLevelMd, restingBPNotHi))
	rule9 := math.Min(ageLevelMd, math.Min(maxHrLevelHi, restingBPNotLo))
	probLevelMd := math.Max(rule4, math.Max(rule5, math.Max(rule6, math.Max(rule7, math.Max(rule8, rule9)))))
	probActivationSetMd := fuzzyLib.NewFuzzySetFromMF(xProb, ConstantMF(probLevelMd)).Intersection(probMdSet)

	rule10 := math.Min(ageLevelLo, math.Min(maxHrLevelLo, restingBPLevelHi))
	rule11 := math.Min(ageLevelMd, math.Min(maxHrLevelLo, restingBPNotLo))
	rule12 := math.Min(ageLevelMd, math.Min(maxHrLevelMd, restingBPLevelHi))
	rule13 := math.Min(ageLevelHi, math.Min(maxHrLevelLo, restingBPNotLo))
	rule14 := math.Min(ageLevelHi, maxHrLevelMd)
	probLevelHi := math.Max(rule10, math.Max(rule11, math.Max(rule12, math.Max(rule13, rule14))))
	probActivationSetHi := fuzzyLib.NewFuzzySetFromMF(xProb, ConstantMF(probLevelHi)).Intersection(probHiSet)

	fmt.Println("Taisyklių sąrašas:")
	fmt.Println("1 taisyklė:", rule1)
	fmt.Println("2 taisyklė:", rule2)
	fmt.Println("3 taisyklė:", rule3)
	fmt.Println("4 taisyklė:", rule4)
	fmt.Println("5 taisyklė:", rule5)
	fmt.Println("6 taisyklė:", rule6)
	fmt.Println("7 taisyklė:", rule7)
	fmt.Println("8 taisyklė:", rule8)
	fmt.Println("9 taisyklė:", rule9)
	fmt.Println("10 taisyklė:", rule10)
	fmt.Println("11 taisyklė:", rule11)
	fmt.Println("12 taisyklė:", rule12)
	fmt.Println("13 taisyklė:", rule13)
	fmt.Println("14 taisyklė:", rule14)

	fmt.Println(probLevelLo, probLevelMd, probLevelHi)

	aggregated := probActivationSetLo.Union(probActivationSetMd).Union(probActivationSetHi)
	fmt.Println(aggregated.Centroid())
}

func arange(start, end, step float64) []float64 {
	var ret []float64
	for i := start; i <= end; i += step {
		ret = append(ret, i)
	}
	return ret
}
