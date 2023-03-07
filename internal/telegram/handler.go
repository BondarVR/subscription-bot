package telegram

import (
	"strings"
	"subscription-bot/internal/models"
	"time"
)

func ParseTime(time string) models.Time {
	arrTime := strings.Split(time, ":")
	structTime := models.Time{
		Hour:    arrTime[0],
		Minutes: arrTime[1],
		Second:  "00",
	}
	return structTime
}

func TimeForTicker() models.Time {
	time := time.Now().Format(layoutOfTime)
	result := ParseTime(time)
	return result
}
