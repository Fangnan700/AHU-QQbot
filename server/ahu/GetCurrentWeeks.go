package ahu

import (
	"Kira-qbot/model"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

var (
	start_date string
)

func init() {
	var config model.Config

	configFileName := "config/Config.yml"
	configFile, _ := os.Open(configFileName)
	decoder := yaml.NewDecoder(configFile)
	_ = decoder.Decode(&config)

	start_date = config.AhuCalendarStartDate
}

func GetCurrentWeek() int {
	layout := "2006-01-02"
	dateStr := time.Now().Format(layout)
	date, _ := time.Parse(layout, dateStr)
	start := time.Date(2023, 3, 6, 0, 0, 0, 0, time.Local)
	week := int(date.Sub(start).Hours()/24/7) + 1

	return week
}
