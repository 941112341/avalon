package inline

import (
	"math"
	"time"
)

// retry >= 0
func RetryFun(f func() error, retry int, wait time.Duration, nextTime func(wt time.Duration, retry int) time.Duration) (err error) {
	if retry < 0 {
		retry = math.MaxInt64 - 1
	}
	for i := 0; i < retry+1; i++ {
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
