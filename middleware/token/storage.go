package token

import "github.com/SevereCloud/vksdk/v2/api"

type Getter interface {
	Get(string, ...api.Params) (string, error)
}

type GetterFunc func(string, ...api.Params) (string, error)

func (g GetterFunc) Get(s string, p ...api.Params) (string, error) {
	return g(s, p...)
}
