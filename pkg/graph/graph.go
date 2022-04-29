package graph

// GetAll

type Graph struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
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
