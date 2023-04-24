package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type Timespan string

const (
	All   Timespan = "4"
	Year  Timespan = "3"
	Month Timespan = "2"
	Day   Timespan = "1"
)

func HashPassword(password string) string {
	hashedPassHex := md5.Sum([]byte(password))
	hashedPassString := hex.EncodeToString(hashedPassHex[:])
	for i := 0; i < len(hashedPassString); i = i + 2 {
		if hashedPassString[i] == 0 {
			hashedPassString = fmt.Sprintf("%sc%s", hashedPassString[0:i], hashedPassString[i+1:])
		}
	}
	return hashedPassString
}

func GetDateString(timespan Timespan, date time.Time) string {
	dateStr := ""
	if timespan == Year {
		dateStr = date.Format("2006")
	} else if timespan == Month {
		dateStr = date.Format("2006-01")
	} else {
		dateStr = date.Format("2006-01-02")
	}
	return dateStr
}

func DateRange(week, year int) (startDate, endDate time.Time) {
	timeBenchmark := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)
	weekStartBenchmark := timeBenchmark.AddDate(0, 0, -(int(timeBenchmark.Weekday())+6)%7)

	_, weekBenchmark := weekStartBenchmark.ISOWeek()
	startDate = weekStartBenchmark.AddDate(0, 0, (week-weekBenchmark)*7)
	endDate = startDate.AddDate(0, 0, 6)

	return startDate, endDate
}

func GetWeekStartEnd(year, week int) (int, int, bool, error) {
	throughMonth := false
	// Week will be calendar week in int
	// Find This Monday
	start, _ := DateRange(week, year)
	mon := start.Day()
	sun := mon + 6

	if sun < mon {
		throughMonth = true
	}

	return mon, sun, throughMonth, nil
}
