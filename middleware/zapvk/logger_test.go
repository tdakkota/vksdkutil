package zapvk

import (
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddleware(t *testing.T) {
	a := require.New(t)
	level := zapcore.InfoLevel
	core, logs := observer.New(level)

	m := Log(zap.New(core), level, true)
	handler := m(func(method string, params ...api.Params) (api.Response, error) {
		return api.Response{}, nil
	})

	_, err := handler("", api.Params{})
	if err != nil {
		t.Fatal(err)
	}

	a.Equal(1, logs.Len())
	log := logs.All()[0]
	a.Equal("VK request", log.Message)
}
