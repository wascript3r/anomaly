package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/anomaly/pkg/graph"
)

type Usecase struct {
	graphRepo  graph.Repository
	ctxTimeout time.Duration

	validate graph.Validate
}

func New(gr graph.Repository, t time.Duration, v graph.Validate) *Usecase {
	return &Usecase{
		graphRepo:  gr,
		ctxTimeout: t,

		validate: v,
	}
}

func (u *Usecase) GetAll(ctx context.Context) (*graph.GetAllRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	gs, err := u.graphRepo.GetAll(c)
	if err != nil {
		return nil, err
	}

	graphs := make([]*graph.Graph, len(gs))
	for i, g := range gs {
		graphs[i] = &graph.Graph{
			ID:      g.ID,
			Name:    g.Name,
			TrapMFs: make([]*graph.TrapMF, len(g.TrapMFs)),
		}
		for j, t := range g.TrapMFs {
			graphs[i].TrapMFs[j] = &graph.TrapMF{
				ID:     t.ID,
				Name:   t.Name,
				Coeffs: t.Coeffs,
			}
		}
	}

	return &graph.GetAllRes{
		Graphs: graphs,
	}, nil
}

func (u *Usecase) UpdateTrapMF(ctx context.Context, req *graph.UpdateTrapMFReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return graph.InvalidInputError
	}

	if !u.validate.TrapMFCoeffs(req.Coeffs) {
		return graph.InvalidCoeffsError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	t, err := u.graphRepo.GetTrapMF(c, req.ID)
	if err != nil {
		return err
	}

	if len(req.Coeffs) != len(t.Coeffs) {
		return graph.InvalidCoeffsDimError
	}

	return u.graphRepo.UpdateTrapMF(c, req.ID, req.Coeffs)
}
