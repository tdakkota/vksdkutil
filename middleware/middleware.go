package middleware

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/tdakkota/vksdkutil"
)

// DeleteParams is middleware which deletes query params from list.
// It's useful for API proxy.
func DeleteParams(list []string) func(handler sdkutil.Handler) sdkutil.Handler {
	if len(list) == 0 {
		return func(handler sdkutil.Handler) sdkutil.Handler {
			return handler
		}
	}

	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params api.Params) (api.Response, error) {
			if len(list) == 1 {
				delete(params, list[0])
			} else {
				for i := range list {
					delete(params, list[i])
				}
			}

			return handler(method, params)
		}
	}
}
