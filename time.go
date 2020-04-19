package etracking

import "time"

const etrackingTimeFormat = "02-01-2006 15:04:05"

func parseETrackingTime(timeString string) (time.Time, error) {
	t, err := time.ParseInLocation(etrackingTimeFormat, timeString,
		time.FixedZone("UTC+7", 3600*7))
	if err != nil {
		return time.Time{}, err
	}

	return beYearToCEYear(t), nil
}

func beYearToCEYear(t time.Time) time.Time {
	return t.AddDate(-543, 0, 0)
}
