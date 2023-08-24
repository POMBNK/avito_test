package utils

import "time"

func Again(fn func() error, attempts int, delay time.Duration) error {
	for attempts > 0 {
		err := fn()
		if err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return nil
}
