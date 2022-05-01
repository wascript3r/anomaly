package graph

// GetAll

type Graph struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	MinVal  *int      `json:"minVal"`
	MaxVal  *int      `json:"maxVal"`
	TrapMFs []*TrapMF `json:"trapMFs"`
}

type TrapMF struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Coeffs []int  `json:"coeffs"`
}

type GetAllRes struct {
	Graphs []*Graph `json:"graphs"`
}

// UpdateTrapMF

type UpdateTrapMFReq struct {
	ID     int   `json:"id" validate:"required"`
	Coeffs []int `json:"coeffs" validate:"required,min=4,max=8"`
}
