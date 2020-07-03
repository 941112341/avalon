package inline

import (
	"errors"
	"time"
)

func RetryFun(f func() error, retry int, wait time.Duration, nextTime func(wt time.Duration, retry int) time.Duration) (err error) {
	if retry == 0 {
		return errors.New("retry cannot be zero")
	}
	for i := 0; i < retry; i++ {
		if err = f(); err == nil {
			return nil
		}
		time.Sleep(wait)
		wait = nextTime(wait, i)
	}
	return
}

func Retry(f func() error, retry int, wait time.Duration) error {
	return RetryFun(f, retry, wait, func(wt time.Duration, _ int) time.Duration {
		return wt << 1
	})
}
