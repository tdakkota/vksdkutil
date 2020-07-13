package token

import "github.com/SevereCloud/vksdk/api"

type Getter interface {
	Get(string, api.Params) (string, error)
}