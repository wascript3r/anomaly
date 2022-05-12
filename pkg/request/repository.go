package request

import (
	"context"

	"github.com/wascript3r/anomaly/pkg/domain"
)

type Repository interface {
	Insert(ctx context.Context, rs *domain.Request) error
	GetStats(ctx context.Context, rs *domain.Request) (*domain.RequestStats, error)
	GetIMSIID(ctx context.Context, imsi string) (int, error)
	GetMSCID(ctx context.Context, msc string) (int, error)
	InsertIMSI(ctx context.Context, is *domain.IMSI) error
	InsertMSC(ctx context.Context, ms *domain.MSC) error
	GetTotalStats(ctx context.Context, anomalyThreshold float64, filter *domain.RequestFilter) ([]*domain.RequestTotalStats, error)
	GetAll(ctx context.Context, filter *domain.RequestFilter) ([]*domain.RequestMeta, error)
}
