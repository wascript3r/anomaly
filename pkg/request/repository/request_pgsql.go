package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/anomaly/pkg/domain"
	"github.com/wascript3r/anomaly/pkg/repository/pgsql"
)

const (
	insertSQL     = "INSERT INTO requests (timestamp, imsi_id, msc_id, anomaly_score) VALUES ($1, $2, $3, $4) RETURNING id"
	getStatsSQL   = `SELECT COUNT(id) FROM requests WHERE ("timestamp" >= $1::TIMESTAMP - INTERVAL '1 hour' AND "timestamp" <= $1::TIMESTAMP) AND imsi_id = $2 UNION ALL SELECT COUNT(id) FROM requests WHERE ("timestamp" >= $1::TIMESTAMP - INTERVAL '1 hour' AND "timestamp" <= $1::TIMESTAMP) AND msc_id = $3`
	getIMSIIDSQL  = "SELECT id FROM imsi WHERE imsi = $1"
	getMSCIDSQL   = "SELECT id FROM msc WHERE msc = $1"
	insertIMSISQL = "INSERT INTO imsi (imsi) VALUES ($1) RETURNING id"
	insertMSCSQL  = "INSERT INTO msc (msc) VALUES ($1) RETURNING id"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) Insert(ctx context.Context, rs *domain.Request) error {
	return p.conn.QueryRowContext(ctx, insertSQL, rs.Timestamp, rs.IMSIID, rs.MSCID, rs.AnomalyScore).Scan(&rs.ID)
}

func (p *PgRepo) GetStats(ctx context.Context, rs *domain.Request) (*domain.RequestStats, error) {
	stats := &domain.RequestStats{}
	rows, err := p.conn.QueryContext(ctx, getStatsSQL, rs.Timestamp, rs.IMSIID, rs.MSCID)
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

func (p *PgRepo) GetIMSIID(ctx context.Context, imsi string) (int, error) {
	var id int
	err := p.conn.QueryRowContext(ctx, getIMSIIDSQL, imsi).Scan(&id)
	if err != nil {
		return 0, pgsql.ParseSQLError(err)
	}
	return id, nil
}

func (p *PgRepo) GetMSCID(ctx context.Context, msc string) (int, error) {
	var id int
	err := p.conn.QueryRowContext(ctx, getMSCIDSQL, msc).Scan(&id)
	if err != nil {
		return 0, pgsql.ParseSQLError(err)
	}
	return id, nil
}

func (p *PgRepo) InsertIMSI(ctx context.Context, is *domain.IMSI) error {
	return p.conn.QueryRowContext(ctx, insertIMSISQL, is.IMSI).Scan(&is.ID)
}

func (p *PgRepo) InsertMSC(ctx context.Context, ms *domain.MSC) error {
	return p.conn.QueryRowContext(ctx, insertMSCSQL, ms.MSC).Scan(&ms.ID)
}
