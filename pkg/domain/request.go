package domain

import "time"

type Request struct {
	ID           int
	Timestamp    time.Time
	IMSIID       int
	MSCID        int
	AnomalyScore float64
}

type RequestStats struct {
	IMSIReqs int
	MSCReqs  int
}
