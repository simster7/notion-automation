package common

import (
	"log"
	"time"
)

type NotionTime struct {
	time time.Time
}

func newNotionTime() NotionTime {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatal(err)
	}
	return NotionTime{time: time.Now().In(loc)}
}

func (e NotionTime) isInit() bool {
	return !e.time.IsZero()
}

func (e NotionTime) NotionDate() string {
	return e.time.Format("2006-01-02")
}

func (e NotionTime) Format(format string) string {
	return e.time.Format(format)
}

func (e NotionTime) AddDate(year, month, day int) NotionTime {
	return NotionTime{time: e.time.AddDate(year, month, day)}
}

var execTime NotionTime

func SetTime() {
	if !execTime.isInit() {
		execTime = newNotionTime()
	}
}

func GetTime() NotionTime {
	return execTime
}
