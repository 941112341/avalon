package inline

import "time"

func Between(t, start, end time.Time) bool {
	return t.After(start) && t.Before(end)
}
