package model

type Course struct {
	Index       int      `json:"index"`
	Class       []string `json:"class"`
	CName       string   `json:"c_name"`
	TName       string   `json:"t_name"`
	CDayNum     int      `json:"c_day_num"`
	CDay        string   `json:"c_day"`
	CTime       string   `json:"c_time"`
	Address     string   `json:"address"`
	StartWeek   int      `json:"start_week"`
	EndWeek     int      `json:"end_week"`
	CurrentWeek int      `json:"current_week"`
}
