package cache

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockStorage struct {
	reads  int
	writes int
}

func (m *mockStorage) Save(key Key, value api.Response) error {
	m.writes++
	return nil
}

func (m *mockStorage) Load(key Key) (api.Response, error) {
	if m.writes < 1 {
		return api.Response{}, ErrCacheMiss
	}

	m.reads++
	return api.Response{}, nil
}

func TestMiddleware(t *testing.T) {
	t.Run("all-cache", func(t *testing.T) {
		storage := new(mockStorage)
		m := Create(storage, func(method string, param api.Params) bool {
			return true
		})

		handler := m(func(method string, params api.Params) (api.Response, error) {
			return api.Response{}, nil
		})

		_, err := handler("test", api.Params{"a": "1", "b": "2"})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, storage.writes)

		_, err = handler("test", api.Params{"b": "2", "a": "1"})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, storage.reads)
	})

	t.Run("do-not-cache", func(t *testing.T) {
		storage := new(mockStorage)
		m := Create(storage, func(method string, param api.Params) bool {
			return false
		})

		handler := m(func(method string, params api.Params) (api.Response, error) {
			return api.Response{}, nil
		})

		_, err := handler("test", api.Params{"a": "1", "b": "2"})
		if err != nil {
			t.Fatal(err)
		}
		assert.Zero(t, storage.writes)

		_, err = handler("test", api.Params{"b": "2", "a": "1"})
		if err != nil {
			t.Fatal(err)
		}
		assert.Zero(t, storage.reads)
	})
}
