package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/wascript3r/anomaly/pkg/domain"
)

const (
	getAllSQL       = "SELECT g.id, g.name, g.infinite, t.id, t.name, t.coeffs FROM graphs g INNER JOIN trap_mfs t ON t.graph_id = g.id ORDER BY g.id, t.id ASC"
	getTrapMFSQL    = "SELECT t.id, t.name, t.coeffs, g.min_val, g.max_val FROM trap_mfs t INNER JOIN graphs g ON g.id = t.graph_id WHERE t.id = $1"
	updateTrapMFSQL = "UPDATE trap_mfs SET coeffs = $2 WHERE id = $1"
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

func (p *PgRepo) GetTrapMF(ctx context.Context, id int) (*domain.FullTrapMF, error) {
	row := p.conn.QueryRowContext(ctx, getTrapMFSQL, id)
	var (
		coeffs []int64
		t      domain.FullTrapMF
	)
	err := row.Scan(&t.ID, &t.Name, pq.Array(&coeffs), &t.MinVal, &t.MaxVal)
	if err != nil {
		return nil, err
	}

	t.Coeffs = make([]int, len(coeffs))
	for i, v := range coeffs {
		t.Coeffs[i] = int(v)
	}

	return &t, nil
}

func (p *PgRepo) UpdateTrapMF(ctx context.Context, id int, coeffs []int) error {
	_, err := p.conn.ExecContext(ctx, updateTrapMFSQL, id, pq.Array(coeffs))
	return err
}
