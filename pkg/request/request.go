package request

// Process

type ProcessReq struct {
	Timestamp string `json:"timestamp" validate:"required,datetime"`
	IMSI      string `json:"imsi" validate:"required,r_imsi"`
	MSC       string `json:"msc" validate:"required,r_msc"`
}

type ProcessRes struct {
	AnomalyScore float64 `json:"anomalyScore"`
}
