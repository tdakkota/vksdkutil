package zerolog

import (
	"time"

	"github.com/SevereCloud/vksdk/api"
	"github.com/rs/zerolog"
	sdkutil "github.com/tdakkota/vksdkutil"
)

func LoggingMiddleware(l zerolog.Logger) func(handler sdkutil.Handler) sdkutil.Handler {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params api.Params) (api.Response, error) {
			start := time.Now()
			r, err := handler(method, params)

			l.WithLevel(l.GetLevel()).
				Str("method", method).
				Interface("params", params).
				Dur("took", time.Since(start)).
				Err(err).
				Msg("send VK request")
			return r, err
		}
	}
}
