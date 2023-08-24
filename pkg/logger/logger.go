package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range h.Writer {
		w.Write([]byte(line))
	}

	return err
}

func (h *writerHook) Levels() []logrus.Level {
	return h.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{e}
}

func (l *Logger) GetLoggerWithFields(k string, v any) Logger {
	return Logger{l.WithField(k, v)}
}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d:", fileName, frame.Line)
		},
		DisableColors: true,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	alllogs, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{alllogs, os.Stdout},
		LogLevels: logrus.AllLevels})

	l.SetLevel(logrus.TraceLevel)
	e = logrus.NewEntry(l)
}
