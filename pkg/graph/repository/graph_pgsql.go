package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/wascript3r/anomaly/pkg/domain"
)

const (
	getAllSQL = "SELECT g.id, g.name, g.infinite, t.id, t.name, t.coeffs FROM graphs g INNER JOIN trap_mfs t ON t.graph_id = g.id"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) GetAll(ctx context.Context) ([]*domain.Graph, error) {
	rows, err := p.conn.QueryContext(ctx, getAllSQL)
	if err != nil {
		return nil, err
	}

	var graphs []*domain.Graph
	for rows.Next() {
		var (
			coeffs []int64
			g      domain.Graph
			t      domain.TrapMF
		)
		err := rows.Scan(&g.ID, &g.Name, &g.Infinite, &t.ID, &t.Name, pq.Array(&coeffs))
		if err != nil {
			return nil, err
		}

		t.Coeffs = make([]int, len(coeffs))
		for i, v := range coeffs {
			t.Coeffs[i] = int(v)
		}

		if len(graphs) == 0 || graphs[len(graphs)-1].ID != g.ID {
			g.TrapMFs = []*domain.TrapMF{&t}
			graphs = append(graphs, &g)
		} else {
			ind := len(graphs) - 1
			graphs[ind].TrapMFs = append(graphs[ind].TrapMFs, &t)
		}
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return graphs, nil
}
