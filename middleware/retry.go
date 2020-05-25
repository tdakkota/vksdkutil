package middleware

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/tdakkota/vksdkutil"
	"time"
)

func Retry(maxAttempts int, timeout time.Duration) func(handler sdkutil.Handler) sdkutil.Handler {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params api.Params) (r api.Response, err error) {
			attempt := 0

			for attempt < maxAttempts {
				r, err = handler(method, params)
				if err != nil {
					attempt++
					time.Sleep(timeout)

					continue
				}
			}

			return r, err
		}
	}
}
