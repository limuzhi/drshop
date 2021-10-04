package global

import "time"

const TimeLayout = "2006-01-02 15:04:05"

func GetDateByUnix(t int64) string {
	if t <= 0 {
		return ""
	}
	return time.Unix(t, 0).Format(TimeLayout)
}
