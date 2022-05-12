package request

import "time"

// Process

type ProcessReq struct {
	Timestamp string `json:"timestamp" validate:"required,datetime"`
	IMSI      string `json:"imsi" validate:"required,r_imsi"`
	MSC       string `json:"msc" validate:"required,r_msc"`
}

type ProcessRes struct {
	AnomalyScore float64 `json:"anomalyScore"`
}

// GetStats

type FilterReq struct {
	StartTime *string `json:"startTime" validate:"omitempty,datetime"`
	EndTime   *string `json:"endTime" validate:"omitempty,datetime"`
	IMSIS     []int   `json:"imsis" validate:"omitempty,gt=0"`
	MSCS      []int   `json:"mscs" validate:"omitempty,gt=0"`
}

type TotalStats struct {
	Timestamp time.Time `json:"timestamp"`
	Total     int       `json:"total"`
	Anomalies int       `json:"anomalies"`
}

type GetStatsRes struct {
	TotalStats []*TotalStats `json:"totalStats"`
}

// GetIMSIStats

type AdvancedFilterReq struct {
	FilterReq
	Limit int `json:"limit" validate:"required,gt=0"`
}

type IMSIStats struct {
	ID        int    `json:"id"`
	IMSI      string `json:"imsi"`
	Total     int    `json:"total"`
	Anomalies int    `json:"anomalies"`
}

type GetIMSIStatsRes struct {
	Count     int          `json:"count"`
	IMSIStats []*IMSIStats `json:"imsiStats"`
}

// GetMSCStats

type MSCStats struct {
	ID        int    `json:"id"`
	MSC       string `json:"msc"`
	Total     int    `json:"total"`
	Anomalies int    `json:"anomalies"`
}

type GetMSCStatsRes struct {
	Count    int         `json:"count"`
	MSCStats []*MSCStats `json:"mscStats"`
}

// GetAll

type Request struct {
	Timestamp    time.Time `json:"timestamp"`
	IMSI         string    `json:"imsi"`
	MSC          string    `json:"msc"`
	AnomalyScore float64   `json:"anomalyScore"`
}

type GetAllRes struct {
	Count    int        `json:"count"`
	Requests []*Request `json:"requests"`
}
