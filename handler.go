package sdkutil

import "github.com/SevereCloud/vksdk/api"

type Handler = func(method string, params api.Params) (api.Response, error)

type Middleware func(handler Handler) Handler

func PatchHandler(sdk *api.VK, f Middleware) *api.VK {
	tmp := sdk.Handler
	sdk.Handler = func(method string, params api.Params) (api.Response, error) {
		return f(tmp)(method, params)
	}

	return sdk
}
