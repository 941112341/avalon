package log

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestNewLoggerWithRotate(t *testing.T) {
	NewLoggerWithRotate().WithFields(logrus.Fields{
		"hello": "world",
	}).Error("message")
}
