package request

import (
	"context"

	"github.com/wascript3r/anomaly/pkg/domain"
)

type Repository interface {
	Insert(ctx context.Context, rs *domain.Request) error
	GetStats(ctx context.Context, imsi, msc string) (*domain.RequestStats, error)
}
