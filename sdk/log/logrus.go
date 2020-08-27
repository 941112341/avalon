package log

import (
	"github.com/getsentry/sentry-go"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var Log *logrus.Logger
var FileLog *logrus.Logger

var once sync.Once

func File() *logrus.Logger {
	if Log != nil {
		return Log
	}
	once.Do(func() {
		path := "/tmp/avalon.log"
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

		logrus.SetLevel(logrus.InfoLevel)

		// sentry 激活暂时不写到配置文件里 等待重构
		if os.Getenv("env") == "online" {
			err := sentry.Init(sentry.ClientOptions{
				// Either set your DSN here or set the SENTRY_DSN environment variable.
				Dsn: "https://451f8f2235654834af86d61b87660870@o439750.ingest.sentry.io/5407022",
				// Enable printing of SDK debug messages.
				// Useful when getting started or trying to figure something out.
				Debug: false,
			})
			if err != nil {
				Log.Errorln(err)
				return
			}

		}
	})

	return Log
}

func New() *logrus.Logger {
	//return Console()
	return File()
}

func Console() *logrus.Logger {
	if Log != nil {
		return Log
	}
	once.Do(func() {
		Log = logrus.New()
		Log.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: true,
		})
		Log.SetLevel(logrus.InfoLevel)
		Log.SetOutput(os.Stdout)
	})

	return Log
}
