package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/anomaly/pkg/domain"
	"github.com/wascript3r/anomaly/pkg/fuzzy"
	"github.com/wascript3r/anomaly/pkg/request"
)

type Usecase struct {
	requestRepo request.Repository
	ctxTimeout  time.Duration

	anomalyThreshold float64

	fuzzyUcase fuzzy.UseCase
	validate   request.Validate
}

func New(rr request.Repository, t time.Duration, anomalyThreshold float64, fu fuzzy.UseCase, v request.Validate) *Usecase {
	return &Usecase{
		requestRepo: rr,
		ctxTimeout:  t,

		anomalyThreshold: anomalyThreshold,

		fuzzyUcase: fu,
		validate:   v,
	}
}

func (u *Usecase) Process(ctx context.Context, req *request.ProcessReq) (*request.ProcessRes, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return nil, request.InvalidInputError
	}

	timestamp, err := time.Parse(u.validate.GetDateTimeFormat(), req.Timestamp)
	if err != nil {
		return nil, request.CannotParseTimestampError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	imsi := &domain.IMSI{
		IMSI: req.IMSI,
	}
	imsi.ID, err = u.requestRepo.GetIMSIID(c, imsi.IMSI)
	if err != nil {
		if err != domain.ErrNotFound {
			return nil, err
		}
		err = u.requestRepo.InsertIMSI(c, imsi)
		if err != nil {
			return nil, err
		}
	}

	msc := &domain.MSC{
		MSC: req.MSC,
	}
	msc.ID, err = u.requestRepo.GetMSCID(c, msc.MSC)
	if err != nil {
		if err != domain.ErrNotFound {
			return nil, err
		}
		err = u.requestRepo.InsertMSC(c, msc)
		if err != nil {
			return nil, err
		}
	}

	r := &domain.Request{
		Timestamp: timestamp,
		IMSIID:    imsi.ID,
		MSCID:     msc.ID,
	}

	stats, err := u.requestRepo.GetStats(c, r)
	if err != nil {
		return nil, err
	}

	weekDay := timestamp.Weekday()
	if weekDay == 0 {
		weekDay = 7
	}

	fuzzyModel := &fuzzy.Model{
		DayTime:   float64(timestamp.Hour()),
		WeekDay:   float64(weekDay),
		IMSICalls: float64(stats.IMSIReqs),
		MSCCalls:  float64(stats.MSCReqs),
	}

	r.AnomalyScore, err = u.fuzzyUcase.CalcResult(c, fuzzyModel)
	if err != nil {
		return nil, err
	}
	r.AnomalyScore = float64(int(r.AnomalyScore*10000)) / 10000

	err = u.requestRepo.Insert(c, r)
	if err != nil {
		return nil, err
	}

	return &request.ProcessRes{
		AnomalyScore: r.AnomalyScore,
	}, nil
}

func (u *Usecase) parseFilter(req *request.FilterReq) (*domain.RequestFilter, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return nil, request.InvalidInputError
	}

	filter := &domain.RequestFilter{
		StartTime: nil,
		EndTime:   nil,
		IMSIIDs:   req.IMSIS,
		MSCIDs:    req.MSCS,
	}

	if req.StartTime != nil {
		startTime, err := time.Parse(u.validate.GetDateTimeFormat(), *req.StartTime)
		if err != nil {
			return nil, request.CannotParseTimestampError
		}
		filter.StartTime = &startTime
	}
	if req.EndTime != nil {
		endTime, err := time.Parse(u.validate.GetDateTimeFormat(), *req.EndTime)
		if err != nil {
			return nil, request.CannotParseTimestampError
		}
		filter.EndTime = &endTime
	}

	return filter, nil
}

func (u *Usecase) GetStats(ctx context.Context, req *request.FilterReq) (*request.GetStatsRes, error) {
	filter, err := u.parseFilter(req)
	if err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ts, err := u.requestRepo.GetTotalStats(c, u.anomalyThreshold, filter)
	if err != nil {
		return nil, err
	}

	stats := make([]*request.TotalStats, len(ts))
	for i, t := range ts {
		stats[i] = (*request.TotalStats)(t)
	}

	return &request.GetStatsRes{
		TotalStats: stats,
	}, nil
}

func (u *Usecase) GetAll(ctx context.Context, req *request.FilterReq) (*request.GetAllRes, error) {
	filter, err := u.parseFilter(req)
	if err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	rs, err := u.requestRepo.GetAll(c, filter)
	if err != nil {
		return nil, err
	}

	requests := make([]*request.Request, len(rs))
	for i, r := range rs {
		requests[i] = (*request.Request)(r)
	}

	return &request.GetAllRes{
		Count:    len(requests),
		Requests: requests,
	}, nil
}
