package cache

import (
	"errors"

	"github.com/SevereCloud/vksdk/v2/api"
	sdkutil "github.com/tdakkota/vksdkutil/v2"
)

var ErrCacheMiss = errors.New("cache miss")

// Storage interface represents a cache storage instance.
type Storage interface {
	Save(key Key, value api.Response) error
	Load(key Key) (api.Response, error)
}

type CacheableFunc = func(method string, param ...api.Params) bool

func Create(c Storage, cacheable CacheableFunc) sdkutil.Middleware {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params ...api.Params) (api.Response, error) {
			if !cacheable(method, params...) {
				return handler(method, params...)
			}

			key := NewKey(method, params...)

			r, err := c.Load(key)
			if err != nil && errors.Is(err, ErrCacheMiss) {
				r, err := handler(method, params...)
				if err != nil {
					return r, err
				}

				_ = c.Save(key, r)
			}

			return r, nil
		}
	}
}
