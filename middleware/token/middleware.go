package token

import (
	"github.com/SevereCloud/vksdk/v2/api"
	sdkutil "github.com/tdakkota/vksdkutil/v2"
	"github.com/tdakkota/vksdkutil/v2/middleware/paramsutil"
)

// Create creates new Middleware, which adds token to parameters if force parameter is true
// or if parameter "access_token" not set.
func Create(force bool, getter Getter) sdkutil.Middleware {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params ...api.Params) (api.Response, error) {
			_, ok := paramsutil.Find("access_token", params...)
			if force || !ok {
				token, err := getter.Get(method, params...)
				if err != nil {
					return api.Response{}, err
				}

				params = paramsutil.Append("access_token", token, params...)
			}

			return handler(method, params...)
		}
	}
}
