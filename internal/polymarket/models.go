package polymarket

type MarketsResponse struct {
	Markets []Market `json:"markets"`
}

type Market struct {
	Question string   `json:"question"`
	Slug     string   `json:"slug"`
	Category string   `json:"category"`
	Volume24 float64  `json:"volume24h"`
	OI       float64  `json:"openInterest"`
	Outcomes []string `json:"outcomes"` // если есть
}
