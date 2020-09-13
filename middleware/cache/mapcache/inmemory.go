package mapcache

import (
	"sync"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/tdakkota/vksdkutil/v2/middleware/cache"
)

type Map struct {
	cache map[string]api.Response
	lock  sync.RWMutex
}

func NewMap() *Map {
	return &Map{cache: map[string]api.Response{}}
}

func (m *Map) Save(key cache.Key, value api.Response) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cache[key.String()] = value

	return nil
}

func (m *Map) Load(key cache.Key) (api.Response, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v, ok := m.cache[key.String()]
	if !ok {
		return api.Response{}, cache.ErrCacheMiss
	}

	return v, nil
}
