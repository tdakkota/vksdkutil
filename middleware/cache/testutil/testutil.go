package testutil

import (
	"errors"
	"testing"

	"github.com/SevereCloud/vksdk/api"
	"github.com/stretchr/testify/require"
	"github.com/tdakkota/vksdkutil/middleware/cache"
)

func TestCache(m cache.Storage) func(t *testing.T) {
	return func(t *testing.T) {
		assertions := require.New(t)

		key := cache.NewKey("users.get", api.Params{
			"id": "1",
		})
		value := api.Response{
			Response: []byte("{}"),
		}

		_, err := m.Load(key)
		assertions.Error(err)
		assertions.True(errors.Is(err, cache.ErrCacheMiss))

		err = m.Save(key, value)
		assertions.NoError(err)

		value2, err := m.Load(key)
		assertions.NoError(err)

		assertions.Equal(value.Response, value2.Response)
	}
}
