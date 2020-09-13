package token

import "github.com/SevereCloud/vksdk/v2/api"

type Getter interface {
	Get(string, ...api.Params) (string, error)
}
