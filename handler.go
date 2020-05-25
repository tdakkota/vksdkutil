package sdkutil

import "github.com/SevereCloud/vksdk/api"

type Handler = func(method string, params api.Params) (api.Response, error)

type Middleware func(handler Handler) Handler
