package utils

import (
	"time"
)

func MapToTimestampz(month string, year int) (string, error) {
	mnth, err := time.Parse("January", month)
	if err != nil {
		return "", err
	}
	date := time.Date(year, mnth.Month(), 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	return date, nil
}
