package request

type Music struct {
	Name    string  `json:"name"`
	Intro   string  `json:"intro"`
	Tag     string  `json:"tag"`
	CoverId uint    `json:"coverId"`
	DemoId  uint    `json:"demoId"`
	Tracks  []Track `json:"tracks"`
}

type Track struct {
	Position string  `json:"position"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
	Url      string  `json:"url"`
	Label    string  `json:"label"`
	FileId   uint    `json:"fileId"`
}
