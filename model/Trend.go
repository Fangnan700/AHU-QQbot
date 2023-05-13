package model

type Trend struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Language    string `json:"language"`
	TodayStar   string `json:"today_star"`
}
