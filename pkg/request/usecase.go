package request

import "context"

type Usecase interface {
	Process(ctx context.Context, req *ProcessReq) (*ProcessRes, error)
	GetStats(ctx context.Context, req *FilterReq) (*GetStatsRes, error)
	GetIMSIStats(ctx context.Context, req *AdvancedFilterReq) (*GetIMSIStatsRes, error)
	GetMSCStats(ctx context.Context, req *AdvancedFilterReq) (*GetMSCStatsRes, error)
	GetAll(ctx context.Context, req *FilterReq) (*GetAllRes, error)
}
