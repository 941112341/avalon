package inline

import (
	"strconv"
	"strings"
	"time"
)

func Between(t, start, end time.Time) bool {
	return t.After(start) && t.Before(end)
}

func Parse(dateContext string) time.Duration {
	base := time.Duration(0)
	dateContext = strings.TrimSpace(dateContext)
	if strings.HasSuffix(dateContext, "s") {
		base = time.Second
		dateContext = strings.TrimSuffix(dateContext, "s")
	} else if strings.HasSuffix(dateContext, "m") {
		base = time.Minute
		dateContext = strings.TrimSuffix(dateContext, "m")
	} else if strings.HasSuffix(dateContext, "mill") {
		base = time.Millisecond
		dateContext = strings.TrimSuffix(dateContext, "mill")
	} else if strings.HasSuffix(dateContext, "micr") {
		base = time.Microsecond
		dateContext = strings.TrimSuffix(dateContext, "micr")
	} else if strings.HasSuffix(dateContext, "h") {
		base = time.Hour
		dateContext = strings.TrimSuffix(dateContext, "h")
	}
	dateInt, err := strconv.ParseInt(dateContext, 10, 64)
	if err != nil {
		WithFields("err", err, "param", dateContext).Errorln("parse fail")
	}
	return time.Duration(dateInt) * base
}
