package domain

import "time"

type Request struct {
	ID           int
	Timestamp    time.Time
	IMSIID       int
	MSCID        int
	AnomalyScore float64
}

type RequestMeta struct {
	Timestamp    time.Time
	IMSI         string
	MSC          string
	AnomalyScore float64
}

type RequestStats struct {
	IMSIReqs int
	MSCReqs  int
}

type RequestFilter struct {
	StartTime *time.Time
	EndTime   *time.Time
	IMSIIDs   []int
	MSCIDs    []int
}

type RequestTotalStats struct {
	Timestamp time.Time
	Total     int
	Anomalies int
}

type RequestAdvancedFilter struct {
	RequestFilter
	Limit int
}

type RequestIMSIStats struct {
	ID        int
	IMSI      string
	Total     int
	Anomalies int
}

type RequestMSCStats struct {
	ID        int
	MSC       string
	Total     int
	Anomalies int
}
