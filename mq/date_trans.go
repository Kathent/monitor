package mq

import (
	"time"
)

func TransDate(timestamp int64, format string) string{
	unix := time.Unix(timestamp, 0)
	return unix.Format(format)
}

func TransTime(time time.Time, format string) string{
	return time.Format(format)
}
