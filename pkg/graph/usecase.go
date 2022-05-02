package graph

import "context"

type Usecase interface {
	GetAll(ctx context.Context) (*GetAllRes, error)
	UpdateTrapMF(ctx context.Context, req *UpdateTrapMFReq) error
	GetRuleList(ctx context.Context) (*GetRuleListRes, error)
	UpdateRule(ctx context.Context, req *UpdateRuleReq) error
}
