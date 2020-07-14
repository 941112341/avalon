package log

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestNewLoggerWithRotate(t *testing.T) {
	New().WithFields(logrus.Fields{
		"hello": "world",
	}).Debugln("message")
}
