package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("не указано правило повторения")
	}
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила повторения")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days > 400 {
			return "", errors.New("неверный формат количества дней")
		}
		
		for {
			startDate = startDate.AddDate(0, 0, days)
			if !startDate.Before(now) && !startDate.Equal(now) {
				break
			}
		}
	case "y":
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if !startDate.Before(now) && !startDate.Equal(now) {
				break
			}
		}
	default:
		return "", errors.New("неверный формат правилы повторения")
	}
	return startDate.Format("20060102"), nil
}