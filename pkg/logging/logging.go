package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var LOGPATH = os.Getenv("LOGPATH")

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func New() *Logger {
	return &Logger{e}
}

func init() {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		FullTimestamp: true,
	}

	err := os.MkdirAll(LOGPATH, 0644)
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile(LOGPATH+"log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	logger.SetOutput(io.Discard)

	logger.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	logger.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(logger)
}
