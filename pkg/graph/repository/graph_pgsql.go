package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/wascript3r/anomaly/pkg/domain"
)

const (
	getAllSQL            = "SELECT g.id, g.name, g.infinite, g.min_val, g.max_val, t.id, t.name, t.coeffs FROM graphs g INNER JOIN trap_mfs t ON t.graph_id = g.id ORDER BY g.id, t.id ASC"
	getGraphNamesSQL     = "SELECT name FROM graphs ORDER BY id ASC"
	getTrapMFSQL         = "SELECT t.id, t.name, t.coeffs, g.min_val, g.max_val FROM trap_mfs t INNER JOIN graphs g ON g.id = t.graph_id WHERE t.id = $1"
	getTrapMFsByGraphSQL = "SELECT id, name, coeffs FROM trap_mfs WHERE graph_id = $1 ORDER BY id ASC"
	updateTrapMFSQL      = "UPDATE trap_mfs SET coeffs = $2 WHERE id = $1"
	getRulesAsTextSQL    = "SELECT r.id, t1.name, t2.name, t3.name, t4.name, r.output FROM rules r INNER JOIN trap_mfs t1 ON (t1.id = tf1_id) INNER JOIN trap_mfs t2 ON (t2.id = tf2_id) INNER JOIN trap_mfs t3 ON (t3.id = tf3_id) INNER JOIN trap_mfs t4 ON (t4.id = tf4_id) ORDER BY r.id ASC"
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
		err := rows.Scan(&g.ID, &g.Name, &g.Infinite, &g.MinVal, &g.MaxVal, &t.ID, &t.Name, pq.Array(&coeffs))
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

func (p *PgRepo) GetGraphNames(ctx context.Context) ([]string, error) {
	rows, err := p.conn.QueryContext(ctx, getGraphNamesSQL)
	if err != nil {
		return nil, err
	}

	var names []string
	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return names, nil
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

func (p *PgRepo) GetTrapMFsByGraph(ctx context.Context, graphID int) ([]*domain.FullTrapMF, error) {
	rows, err := p.conn.QueryContext(ctx, getTrapMFsByGraphSQL, graphID)
	if err != nil {
		return nil, err
	}

	var trapMFs []*domain.FullTrapMF
	for rows.Next() {
		var (
			coeffs []int64
			t      domain.FullTrapMF
		)
		err := rows.Scan(&t.ID, &t.Name, pq.Array(&coeffs))
		if err != nil {
			return nil, err
		}

		t.Coeffs = make([]int, len(coeffs))
		for i, v := range coeffs {
			t.Coeffs[i] = int(v)
		}

		trapMFs = append(trapMFs, &t)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return trapMFs, nil
}

func (p *PgRepo) UpdateTrapMF(ctx context.Context, id int, coeffs []int) error {
	_, err := p.conn.ExecContext(ctx, updateTrapMFSQL, id, pq.Array(coeffs))
	return err
}

func (p *PgRepo) GetRulesAsText(ctx context.Context) ([]*domain.RuleText, error) {
	rows, err := p.conn.QueryContext(ctx, getRulesAsTextSQL)
	if err != nil {
		return nil, err
	}

	var rules []*domain.RuleText
	for rows.Next() {
		var r domain.RuleText
		r.Inputs = make([]string, 4)

		err := rows.Scan(&r.ID, &r.Inputs[0], &r.Inputs[1], &r.Inputs[2], &r.Inputs[3], &r.Output)
		if err != nil {
			return nil, err
		}

		rules = append(rules, &r)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return rules, nil
}
