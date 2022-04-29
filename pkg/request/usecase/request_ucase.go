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

	fuzzyUcase fuzzy.UseCase
	validate   request.Validate
}

func New(rr request.Repository, t time.Duration, fu fuzzy.UseCase, v request.Validate) *Usecase {
	return &Usecase{
		requestRepo: rr,
		ctxTimeout:  t,

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

	r.AnomalyScore = u.fuzzyUcase.CalcResult(fuzzyModel)
	r.AnomalyScore = float64(int(r.AnomalyScore*10000)) / 10000

	err = u.requestRepo.Insert(c, r)
	if err != nil {
		return nil, err
	}

	return &request.ProcessRes{
		AnomalyScore: r.AnomalyScore,
	}, nil
}
