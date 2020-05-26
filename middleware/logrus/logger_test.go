package logrus

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

type hookFunc func(*logrus.Entry) error

func (h hookFunc) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
	}
}

func (h hookFunc) Fire(entry *logrus.Entry) error {
	return h(entry)
}

func TestLoggingMiddleware(t *testing.T) {
	logger := logrus.New()
	logger.Hooks.Add(hookFunc(func(entry *logrus.Entry) error {
		assert.Equal(t, logrus.InfoLevel, entry.Level)
		assert.Equal(t, "send VK request", entry.Message)
		return nil
	}))

	m := LoggingMiddleware(logger)
	handler := m(func(method string, params api.Params) (api.Response, error) {
		return api.Response{}, nil
	})

	_, err := handler("", api.Params{})
	if err != nil {
		t.Fatal(err)
	}
}
