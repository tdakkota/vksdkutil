package token

import (
	"github.com/SevereCloud/vksdk/api"
	sdkutil "github.com/tdakkota/vksdkutil"
)

type Getter func(string, api.Params) (string, error)

func Create(force bool, get Getter) sdkutil.Middleware {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params api.Params) (api.Response, error) {
			if force || params["access_token"] == "" {
				token, err := get(method, params)
				if err != nil {
					return api.Response{}, err
				}

				params["access_token"] = token
			}

			return handler(method, params)
		}
	}
}
