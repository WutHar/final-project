package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", fmt.Errorf("некорректный формат даты: %s", dateStr)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила")
	}

	switch parts[0] {
	case "d":
		return handleDailyRule(now, date, parts)
	case "y":
		return handleYearlyRule(now, date)
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", parts[0])
	}
}

func handleDailyRule(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("не указан интервал для правила 'd'")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("некорректный интервал: %s", parts[1])
	}

	if interval <= 0 || interval > 400 {
		return "", fmt.Errorf("интервал должен быть от 1 до 400")
	}

	result := date
	for {
		result = result.AddDate(0, 0, interval)
		if afterNow(result, now) {
			break
		}
	}

	return result.Format("20060102"), nil
}

func handleYearlyRule(now time.Time, date time.Time) (string, error) {
	result := date
	for {
		result = result.AddDate(1, 0, 0)
		if afterNow(result, now) {
			break
		}
	}

	return result.Format("20060102"), nil
}

func afterNow(date, now time.Time) bool {

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}
