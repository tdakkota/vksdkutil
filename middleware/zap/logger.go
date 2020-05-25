package zap

import (
	"github.com/SevereCloud/vksdk/api"
	sdkutil "github.com/tdakkota/vksdkutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func LoggingMiddleware(l *zap.Logger, lvl zapcore.Level) func(handler sdkutil.Handler) sdkutil.Handler {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params api.Params) (api.Response, error) {
			start := time.Now()
			r, err := handler(method, params)

			if e := l.Check(lvl, "send VK request"); e != nil {
				e.Write(
					zap.String("method", method),
					zap.Reflect("params", params),
					zap.Duration("took", time.Since(start)),
					zap.Error(err),
				)
			}

			return r, err
		}
	}
}
