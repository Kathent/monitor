package util

import "time"

func TomorrowDuration(now time.Time) time.Duration {
	tomorrow := time.Now().AddDate(0,0, 1)
	year, month, day := tomorrow.Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, nil)
	return date.Sub(now)
}