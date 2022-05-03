package graph

import (
	"context"

	"github.com/wascript3r/anomaly/pkg/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*domain.Graph, error)
	GetGraphNames(ctx context.Context) ([]string, error)
	GetTrapMF(ctx context.Context, id int) (*domain.FullTrapMF, error)
	GetGraphIDByTrapMF(ctx context.Context, id int) (int, error)
	GetTrapMFsByGraph(ctx context.Context, graphID int) ([]*domain.FullTrapMF, error)
	UpdateTrapMF(ctx context.Context, id int, coeffs []int) error
	GetRulesAsText(ctx context.Context) ([]*domain.RuleText, error)
	GetRules(ctx context.Context) ([]*domain.Rule, error)
	RuleExists(ctx context.Context, id int) (bool, error)
	UpdateRuleOutput(ctx context.Context, id, output int) error
}
