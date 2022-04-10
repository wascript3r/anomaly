package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/anomaly/pkg/domain"
)

const (
	insertSQL = "INSERT INTO requests (id, timestamp, imsi, msc) VALUES ($1, $2, $3, $4) RETURNING id"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) Insert(ctx context.Context, rs *domain.Request) error {
	return p.conn.QueryRowContext(ctx, insertSQL, rs.ID, rs.Timestamp, rs.IMSI, rs.MSC).Scan(&rs.ID)
}
