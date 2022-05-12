package request

import "context"

type Usecase interface {
	Process(ctx context.Context, req *ProcessReq) (*ProcessRes, error)
	GetStats(ctx context.Context, req *FilterReq) (*GetStatsRes, error)
	GetAll(ctx context.Context, req *FilterReq) (*GetAllRes, error)
}
