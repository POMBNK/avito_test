package utils

import (
	"time"
)

//2023-08-01 00:00:00 +0000 +0000
//
//2023-08-25 15:49:07.758484 +00:00

//func MapToTimestampz(month string, year int) (string, error) {
//	layout := "2006-01-02 15:04:05.999999999-07:00"
//	val, err := time.Parse("January", month)
//	if err != nil {
//		return "", err
//	}
//
//	dateString := fmt.Sprintf("%d-%02d-01 00:00:00.000000000+00:00", year, val.Month())
//	t, err := time.Parse(layout, dateString)
//	if err != nil {
//		return "", err
//	}
//	return t.String(), nil
//}

func MapToTimestampz(month string, year int) (string, error) {
	mnth, err := time.Parse("January", month)
	if err != nil {
		return "", err
	}
	date := time.Date(year, mnth.Month(), 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	return date, nil
}
