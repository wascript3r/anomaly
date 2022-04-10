package request

import "context"

type Usecase interface {
	Process(ctx context.Context, req *ProcessReq) (*ProcessRes, error)
}
