package domain

import "time"

type Request struct {
	ID        int
	Timestamp time.Time
	IMSI      string
	MSC       string
}
