package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
)

var Log *logrus.Logger

func NewLoggerWithRotate() *logrus.Logger {
	if Log != nil {
		return Log
	}

	path := "./avalon.log"
	writer, _ := rotatelogs.New(
		path+".%Y%m%d",
		//path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),             // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(30*24*time.Hour),    // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	pathMap := lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}

	Log = logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	return Log
}