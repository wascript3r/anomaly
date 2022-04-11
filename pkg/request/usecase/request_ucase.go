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
	fuzzyUcase       fuzzy.UseCase
	validate         request.Validate
}

func New(rr request.Repository, t time.Duration, anomalyThreshold float64, fu fuzzy.UseCase, v request.Validate) *Usecase {
	return &Usecase{
		requestRepo: rr,
		ctxTimeout:  t,

		anomalyThreshold: anomalyThreshold,
		fuzzyUcase:       fu,
		validate:         v,
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

	r := &domain.Request{
		Timestamp: timestamp,
		IMSI:      req.IMSI,
		MSC:       req.MSC,
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
	anomalyScore := u.fuzzyUcase.CalcResult(fuzzyModel)
	blocked := anomalyScore >= u.anomalyThreshold

	if !blocked {
		err = u.requestRepo.Insert(c, r)
		if err != nil {
			return nil, err
		}
	}

	return &request.ProcessRes{
		Blocked:      blocked,
		AnomalyScore: anomalyScore,
	}, nil
}
