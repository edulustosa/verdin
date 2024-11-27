package utils

import "time"

func GetLastMonth() time.Month {
	return time.Now().AddDate(0, -1, 0).Month()
}
