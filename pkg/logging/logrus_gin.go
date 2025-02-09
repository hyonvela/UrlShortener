package logging

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type GinLogrusWriter struct {
	logger *logrus.Entry
}

func (w *GinLogrusWriter) Write(p []byte) (n int, err error) {
	message := strings.TrimSuffix(string(p), "\n")

	switch {
	case strings.Contains(message, "[WARNING]"):
		w.logger.Warn(message)
	case strings.Contains(message, "[ERROR]"):
		w.logger.Error(message)
	case strings.Contains(message, "[DEBUG]"):
		w.logger.Debug(message)
	default:
		w.logger.Info(message)
	}

	return len(p), nil
}

func NewGinLogrusWriter(logger *logrus.Entry) io.Writer {
	return &GinLogrusWriter{logger: logger}
}
