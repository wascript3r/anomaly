package graph

import (
	"context"

	"github.com/wascript3r/anomaly/pkg/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*domain.Graph, error)
	GetGraphNames(ctx context.Context) ([]string, error)
	GetTrapMF(ctx context.Context, id int) (*domain.FullTrapMF, error)
	GetTrapMFsByGraph(ctx context.Context, graphID int) ([]*domain.FullTrapMF, error)
	UpdateTrapMF(ctx context.Context, id int, coeffs []int) error
	GetRulesAsText(ctx context.Context) ([]*domain.RuleText, error)
}
