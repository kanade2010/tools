package main

import (
  	"os"
  	"github.com/sirupsen/logrus"
    "github.com/lestrrat-go/file-rotatelogs"
    "github.com/rifflock/lfshook"
	"time"
	"io"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/ailumiyana/tools/rotate"
)


type RotateFileConfig struct {
	Filename   string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int //days
	Level      logrus.Level
	Formatter  logrus.Formatter
}

type RotateFileHook struct {
	Config    RotateFileConfig
	logWriter io.Writer
}

func NewRotateFileHook(config RotateFileConfig) (logrus.Hook, error) {

	hook := RotateFileHook{
		Config: config,
	}
	hook.logWriter = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	return &hook, nil
}

func (hook *RotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.Config.Level+1]
}

func (hook *RotateFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.logWriter.Write(b)
	return nil
}








func newLfsHook(maxRemainCnt uint) logrus.Hook {
    writer, err := rotatelogs.New(
        "./logrus.log" + ".%Y%m%d%H",
        // WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
        rotatelogs.WithLinkName("./logrus.log"),

        // WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
        rotatelogs.WithRotationTime(time.Hour),

        // WithMaxAge和WithRotationCount二者只能设置一个，
        // WithMaxAge设置文件清理前的最长保存时间，
        // WithRotationCount设置文件清理前最多保存的个数。
        //rotatelogs.WithMaxAge(24*time.Hour),
        rotatelogs.WithRotationCount(maxRemainCnt),
    )

    if err != nil {
        log.Errorf("config local file system for logger error: %v", err)
    }

    lfsHook := lfshook.NewHook(lfshook.WriterMap{
        logrus.DebugLevel: writer,
        logrus.InfoLevel:  writer,
        logrus.WarnLevel:  writer,
        logrus.ErrorLevel: writer,
        logrus.FatalLevel: writer,
        logrus.PanicLevel: writer,
    }, &logrus.TextFormatter{DisableColors: true})

    return lfsHook
}

// Create a new instance of the logger. You can have any number of instances.
var log = logrus.New()

func main() {
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	log.Out = os.Stdout
	/*log.SetFormatter(&logrus.JSONFormatter{})

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")*/

	//log.SetFlags(logrus.Ldate|logrus.Lshortfile)
	//hook := newLfsHook( 5)
	hook, _ = rotate.NewRotateFileHook(rotate.RotateFileConfig{
		Filename: "./logrus.log",
		MaxSize: 1,
		MaxBackups: 7,
		MaxAge: 7,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
			TimestampFormat : time.RFC822,
		},})
	log.AddHook(hook)

	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		TimestampFormat : time.RFC822,
	})

	for {
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Info("A group of walrus emerges from the ocean")
		time.Sleep(10*time.Millisecond)
	}
	// You could set this to any `io.Writer` such as a file
	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	//  log.Out = file
	// } else {
	//  log.Info("Failed to log to file, using default stderr")
	// }


}
