package graph

import (
	"context"

	"github.com/wascript3r/anomaly/pkg/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*domain.Graph, error)
	GetTrapMF(ctx context.Context, id int) (*domain.FullTrapMF, error)
	UpdateTrapMF(ctx context.Context, id int, coeffs []int) error
}
