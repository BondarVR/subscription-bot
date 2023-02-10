package telegram

import (
	"strings"
	"time"
)

func ParseTime(time string) Time {
	arrTime := strings.Split(time, ":")
	structTime := Time{
		Hour:    arrTime[0],
		Minutes: arrTime[1],
		Second:  "00",
	}
	return structTime
}

func TimeForTicker() Time {
	time := time.Now().Format("15:04:05")
	result := ParseTime(time)
	return result
}
