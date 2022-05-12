package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/wascript3r/anomaly/pkg/domain"
	"github.com/wascript3r/anomaly/pkg/repository/pgsql"
)

const (
	insertSQL        = "INSERT INTO requests (timestamp, imsi_id, msc_id, anomaly_score) VALUES ($1, $2, $3, $4) RETURNING id"
	getStatsSQL      = `SELECT COUNT(id) FROM requests WHERE ("timestamp" >= $1::TIMESTAMP - INTERVAL '1 hour' AND "timestamp" <= $1::TIMESTAMP) AND imsi_id = $2 UNION ALL SELECT COUNT(id) FROM requests WHERE ("timestamp" >= $1::TIMESTAMP - INTERVAL '1 hour' AND "timestamp" <= $1::TIMESTAMP) AND msc_id = $3`
	getIMSIIDSQL     = "SELECT id FROM imsi WHERE imsi = $1"
	getMSCIDSQL      = "SELECT id FROM msc WHERE msc = $1"
	insertIMSISQL    = "INSERT INTO imsi (imsi) VALUES ($1) RETURNING id"
	insertMSCSQL     = "INSERT INTO msc (msc) VALUES ($1) RETURNING id"
	getTotalStatsSQL = "SELECT DATE_TRUNC('hour', timestamp) AS timestamp_h, COUNT(id) AS total, COUNT(id) FILTER (WHERE anomaly_score >= $1) AS anomalies FROM requests <filter> GROUP BY timestamp_h ORDER BY timestamp_h ASC"
	getIMSIStatsSQL  = "SELECT i.id, i.imsi, r.total, r.anomalies FROM (SELECT imsi_id, COUNT(id) AS total, COUNT(id) FILTER (WHERE anomaly_score >= $1) AS anomalies FROM requests <filter> GROUP BY imsi_id) r INNER JOIN imsi i ON i.id = r.imsi_id ORDER BY total DESC LIMIT $2"
	getMSCStatsSQL   = "SELECT m.id, m.msc, r.total, r.anomalies FROM (SELECT msc_id, COUNT(id) AS total, COUNT(id) FILTER (WHERE anomaly_score >= $1) AS anomalies FROM requests <filter> GROUP BY msc_id) r INNER JOIN msc m ON m.id = r.msc_id ORDER BY total DESC LIMIT $2"
	getAllSQL        = "SELECT timestamp, imsi, msc, anomaly_score FROM requests r INNER JOIN imsi i ON i.id = r.imsi_id INNER JOIN msc m ON m.id = r.msc_id <filter> ORDER BY timestamp ASC"

	startTimestampFilter = "timestamp >= <param>::TIMESTAMP"
	endTimestampFilter   = "timestamp <= <param>::TIMESTAMP"
	imsiFilter           = "imsi_id = ANY(<param>::INT[])"
	mscFilter            = "msc_id = ANY(<param>::INT[])"
)

type filter struct {
	baseQuery  string
	query      string
	addedWhere bool
	args       []interface{}
}

func newFilter(baseQuery string) *filter {
	return &filter{
		baseQuery:  baseQuery,
		query:      "",
		addedWhere: false,
		args:       make([]interface{}, 0),
	}
}

func (f *filter) add(filterQuery string, value interface{}) {
	if !f.addedWhere {
		f.query += " WHERE "
		f.addedWhere = true
	} else {
		f.query += " AND "
	}

	f.query += strings.Replace(filterQuery, "<param>", fmt.Sprintf("$%d", len(f.args)+1), 1)
	f.addArg(value)
}

func (f *filter) getQuery() string {
	return strings.Replace(f.baseQuery, "<filter>", f.query, 1)
}

func (f *filter) addArg(value interface{}) {
	f.args = append(f.args, value)
}

func (f *filter) getArgs() []interface{} {
	return f.args
}

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

func (p *PgRepo) initRequestFilter(filter *domain.RequestFilter, f *filter) {
	if filter.StartTime != nil {
		f.add(startTimestampFilter, *filter.StartTime)
	}
	if filter.EndTime != nil {
		f.add(endTimestampFilter, *filter.EndTime)
	}
	if filter.IMSIIDs != nil {
		f.add(imsiFilter, pq.Array(filter.IMSIIDs))
	}
	if filter.MSCIDs != nil {
		f.add(mscFilter, pq.Array(filter.MSCIDs))
	}
}

func (p *PgRepo) GetTotalStats(ctx context.Context, anomalyThreshold float64, filter *domain.RequestFilter) ([]*domain.RequestTotalStats, error) {
	f := newFilter(getTotalStatsSQL)
	f.addArg(anomalyThreshold)
	p.initRequestFilter(filter, f)

	rows, err := p.conn.QueryContext(ctx, f.getQuery(), f.getArgs()...)
	if err != nil {
		return nil, err
	}

	var stats []*domain.RequestTotalStats
	for rows.Next() {
		var t domain.RequestTotalStats
		err := rows.Scan(&t.Timestamp, &t.Total, &t.Anomalies)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &t)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (p *PgRepo) GetIMSIStats(ctx context.Context, anomalyThreshold float64, filter *domain.RequestAdvancedFilter) ([]*domain.RequestIMSIStats, error) {
	f := newFilter(getIMSIStatsSQL)
	f.addArg(anomalyThreshold)
	f.addArg(filter.Limit)
	p.initRequestFilter(&filter.RequestFilter, f)

	rows, err := p.conn.QueryContext(ctx, f.getQuery(), f.getArgs()...)
	if err != nil {
		return nil, err
	}

	var stats []*domain.RequestIMSIStats
	for rows.Next() {
		var is domain.RequestIMSIStats
		err := rows.Scan(&is.ID, &is.IMSI, &is.Total, &is.Anomalies)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &is)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (p *PgRepo) GetMSCStats(ctx context.Context, anomalyThreshold float64, filter *domain.RequestAdvancedFilter) ([]*domain.RequestMSCStats, error) {
	f := newFilter(getMSCStatsSQL)
	f.addArg(anomalyThreshold)
	f.addArg(filter.Limit)
	p.initRequestFilter(&filter.RequestFilter, f)

	rows, err := p.conn.QueryContext(ctx, f.getQuery(), f.getArgs()...)
	if err != nil {
		return nil, err
	}

	var stats []*domain.RequestMSCStats
	for rows.Next() {
		var ms domain.RequestMSCStats
		err := rows.Scan(&ms.ID, &ms.MSC, &ms.Total, &ms.Anomalies)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &ms)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (p *PgRepo) GetAll(ctx context.Context, filter *domain.RequestFilter) ([]*domain.RequestMeta, error) {
	f := newFilter(getAllSQL)
	p.initRequestFilter(filter, f)

	rows, err := p.conn.QueryContext(ctx, f.getQuery(), f.getArgs()...)
	if err != nil {
		return nil, err
	}

	var requests []*domain.RequestMeta
	for rows.Next() {
		var r domain.RequestMeta
		err := rows.Scan(&r.Timestamp, &r.IMSI, &r.MSC, &r.AnomalyScore)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &r)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return requests, nil
}
