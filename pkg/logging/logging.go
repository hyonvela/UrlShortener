package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type WriterHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger(format string, lvl string) *Logger {
	setupLoger(format, lvl)
	return &Logger{e}
}

func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func setupLoger(format string, lvl string) {
	l := logrus.New()
	l.SetReportCaller(true)

	switch format {
	case "json":
		l.Formatter = &logrus.JSONFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				filename := path.Base(frame.File)
				return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
			},
		}
	case "text":
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				filename := path.Base(frame.File)
				return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
			},
		}
	default:
		log.Fatal("Configure logs format")
	}

	err := os.Mkdir("logs", 0777)
	if err != nil && err.Error() != "mkdir logs: file exists" {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)

	switch lvl {
	case "all":
		l.AddHook((&WriterHook{
			Writer:    []io.Writer{allFile, os.Stdout},
			LogLevels: logrus.AllLevels,
		}))
	case "test":
		l.AddHook((&WriterHook{
			Writer:    []io.Writer{allFile, os.Stdout},
			LogLevels: []logrus.Level{logrus.FatalLevel},
		}))
	}
	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
