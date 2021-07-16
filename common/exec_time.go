package common

import (
	"log"
	"time"
)

const TimeZone = "America/Los_Angeles"

type NotionTime struct {
	time time.Time
}

func newNotionTime() NotionTime {
	loc, err := time.LoadLocation(TimeZone)
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

func (e NotionTime) GetCalendarEventTimes(index int) (string, string) {
	t := e.time
	endIndex := index + 1
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 8+index/2, 30*(index%2), 0, 0, t.Location())
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 8+endIndex/2, 30*(endIndex%2), 0, 0, t.Location())
	return startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)
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
