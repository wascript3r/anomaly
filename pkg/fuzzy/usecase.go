package fuzzy

import "context"

type UseCase interface {
	CalcResult(ctx context.Context, m *Model) (float64, error)
}
