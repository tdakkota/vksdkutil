package zerolog

import (
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	logger := log.Logger.
		Level(zerolog.InfoLevel)

	logger.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		assert.Equal(t, zerolog.InfoLevel, level)
		assert.Equal(t, "send VK request", message)
	}))

	m := LoggingMiddleware(logger)
	handler := m(func(method string, params ...api.Params) (api.Response, error) {
		return api.Response{}, nil
	})

	_, err := handler("", api.Params{})
	if err != nil {
		t.Fatal(err)
	}
}
