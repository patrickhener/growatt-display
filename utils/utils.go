package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
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

func GetWeekStartEnd(year, week int) (time.Time, time.Time, bool, string, error) {
	var curPrev string
	throughMonth := false
	currentMonth := time.Now().Month()

	// Week will be calendar week in int
	// Find This Monday
	start, _ := DateRange(week, year)
	end := start.AddDate(0, 0, 6)

	if start.Month() == end.Month() {
		if start.Month() == currentMonth {
			curPrev = "cur"
		} else {
			curPrev = "prev"
		}
	} else {
		throughMonth = true
	}

	return start, end, throughMonth, curPrev, nil
}

func DaysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func AddKilowatts(start, end int, data map[string]string) (float64, error) {
	var weekEnergy float64

	for d, e := range data {
		dInt, err := strconv.Atoi(d)
		if err != nil {
			return 0, err
		}
		eFloat, err := strconv.ParseFloat(e, 32)
		if err != nil {
			return 0, err
		}

		// Conditonally calc to Week Energy
		for i := start; i <= end; i++ {
			if dInt == i {
				weekEnergy = weekEnergy + eFloat
			}
		}
	}

	return weekEnergy, nil
}
