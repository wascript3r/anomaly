package graph

import (
	"context"

	"github.com/wascript3r/anomaly/pkg/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*domain.Graph, error)
}
