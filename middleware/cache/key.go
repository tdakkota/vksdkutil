package cache

import (
	"strings"

	"github.com/SevereCloud/vksdk/api"
)

// Key is cache key
type Key struct {
	Method string
	Params OrderedParams
}

func (k Key) String() string {
	b := new(strings.Builder)
	b.WriteString(k.Method)

	for _, v := range k.Params {
		b.WriteString(v.Key)
		b.WriteString(api.FmtValue(v.Value, 0))
	}

	return b.String()
}

func NewKey(method string, params api.Params) Key {
	return Key{Method: method, Params: CreateParams(params)}
}
