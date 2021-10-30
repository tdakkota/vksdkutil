package zapvk

import (
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggingMiddleware(t *testing.T) {
	level := zapcore.InfoLevel

	logger, _ := zap.NewProduction(zap.Hooks(func(entry zapcore.Entry) error {
		assert.Equal(t, level, entry.Level)
		assert.Equal(t, "send VK request", entry.Message)
		return nil
	}))

	m := Log(logger, zap.InfoLevel)
	handler := m(func(method string, params ...api.Params) (api.Response, error) {
		return api.Response{}, nil
	})

	_, err := handler("", api.Params{})
	if err != nil {
		t.Fatal(err)
	}
}
