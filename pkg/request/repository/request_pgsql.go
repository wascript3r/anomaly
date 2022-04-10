package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/anomaly/pkg/domain"
)

const (
	insertSQL   = "INSERT INTO requests (timestamp, imsi, msc) VALUES ($1, $2, $3) RETURNING id"
	getStatsSQL = "SELECT COUNT(id) FROM requests WHERE timestamp >= NOW() - INTERVAL '1 hour' AND imsi = $1 UNION ALL SELECT COUNT(id) FROM requests WHERE timestamp >= NOW() - INTERVAL '1 hour' AND msc = $2"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) Insert(ctx context.Context, rs *domain.Request) error {
	return p.conn.QueryRowContext(ctx, insertSQL, rs.Timestamp, rs.IMSI, rs.MSC).Scan(&rs.ID)
}

func (p *PgRepo) GetStats(ctx context.Context, imsi, msc string) (*domain.RequestStats, error) {
	stats := &domain.RequestStats{}
	rows, err := p.conn.QueryContext(ctx, getStatsSQL, imsi, msc)
	if err != nil {
		return nil, err
	}

	rows.Next()
	if err := rows.Scan(&stats.IMSIReqs); err != nil {
		return nil, err
	}
	rows.Next()
	if err := rows.Scan(&stats.MSCReqs); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return stats, nil
}
