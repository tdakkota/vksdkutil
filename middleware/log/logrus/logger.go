package logrus

import (
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/sirupsen/logrus"
	sdkutil "github.com/tdakkota/vksdkutil/v2"
)

// LoggingMiddleware is middleware which logs VK API request info.
func LoggingMiddleware(l *logrus.Logger) func(handler sdkutil.Handler) sdkutil.Handler {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params ...api.Params) (api.Response, error) {
			start := time.Now()
			r, err := handler(method, params...)

			l.WithFields(logrus.Fields{
				"method": method,
				"params": params,
				"took":   time.Since(start),
				"error":  err,
			}).Log(l.Level, "send VK request")

			return r, err
		}
	}
}
