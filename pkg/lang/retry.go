package lang

import (
	"time"
)

type Retriable func() (bool, error)

func Retry(maxAttempts int, interval time.Duration, fn Retriable) error {
	var (
		forceReturn bool
		err         error
	)
	for i := 0; i < maxAttempts; i++ {
		forceReturn, err = fn()
		if forceReturn || err == nil {
			return err
		}
		sleepInterval := time.Duration(int64(interval) * int64(i+1))
		time.Sleep(sleepInterval)
	}
	return err
}
