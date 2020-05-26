package botpoll

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/SevereCloud/vksdk/object"
)

type subscriptions struct {
	subscriptions map[string]*subscription
	lock          sync.Mutex
}

func newSubscriptions() *subscriptions {
	return &subscriptions{subscriptions: map[string]*subscription{}}
}

func generateKey(check func(string) bool) string {
	for {
		key := strconv.FormatInt(rand.Int63(), 16) + strconv.FormatInt(time.Now().UnixNano(), 16)
		if !check(key) {
			return key
		}
	}
}

func (s *subscriptions) Create() string {
	s.lock.Lock()
	defer s.lock.Unlock()

	key := generateKey(func(k string) bool {
		_, ok := s.subscriptions[k]
		return ok
	})

	s.subscriptions[key] = newSubscription()

	return key
}

func (s *subscriptions) Get(key string) (*subscription, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	v, ok := s.subscriptions[key]
	return v, ok
}

func (s *subscriptions) Notify(events []object.GroupEvent) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, sub := range s.subscriptions {
		sub.Notify(events)
	}
}

func (s *subscriptions) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.subscriptions, key)
}

func (s *subscriptions) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, sub := range s.subscriptions {
		_ = sub.Close()
	}
	return nil
}
