package middleware

import (
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	sdkutil "github.com/tdakkota/vksdkutil/v3"
)

// Retry is middleware which retries requests if underlying handler
// fails.
func Retry(maxAttempts int, timeout time.Duration) func(handler sdkutil.Handler) sdkutil.Handler {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params ...api.Params) (r api.Response, err error) {
			attempt := 0

			for attempt < maxAttempts {
				r, err = handler(method, params...)
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
