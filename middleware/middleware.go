package middleware

import (
	"github.com/SevereCloud/vksdk/v2/api"
	sdkutil "github.com/tdakkota/vksdkutil/v3"
	"github.com/tdakkota/vksdkutil/v3/middleware/paramsutil"
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
		return func(method string, params ...api.Params) (api.Response, error) {
			if len(list) == 1 {
				paramsutil.Delete(list[0], params...)
			} else {
				for i := range list {
					paramsutil.Delete(list[i], params...)
				}
			}

			return handler(method, params...)
		}
	}
}
