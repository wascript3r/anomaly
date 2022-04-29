package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/anomaly/pkg/graph"
)

type Usecase struct {
	graphRepo  graph.Repository
	ctxTimeout time.Duration
}

func New(gr graph.Repository, t time.Duration) *Usecase {
	return &Usecase{
		graphRepo:  gr,
		ctxTimeout: t,
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
