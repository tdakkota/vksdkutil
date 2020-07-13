package token

import (
	"github.com/SevereCloud/vksdk/api"
	sdkutil "github.com/tdakkota/vksdkutil"
)

func Create(force bool, getter Getter) sdkutil.Middleware {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params api.Params) (api.Response, error) {
			if force || params["access_token"] == "" {
				token, err := getter.Get(method, params)
				if err != nil {
					return api.Response{}, err
				}

				params["access_token"] = token
			}

			return handler(method, params)
		}
	}
}
