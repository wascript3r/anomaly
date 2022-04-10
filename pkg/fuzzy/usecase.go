package fuzzy

type UseCase interface {
	CalcResult(m *Model) float64
}
