package botpoll

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/SevereCloud/vksdk/object"
)

type subscription struct {
	notify chan struct{}

	events []object.GroupEvent
	lock   sync.Mutex
}

func newSubscription() *subscription {
	return &subscription{notify: make(chan struct{}, 1)}
}

func (s *subscription) Notify(events []object.GroupEvent) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.events = append(s.events, events...)
	s.notify <- struct{}{}
}

func pollResponse(events []object.GroupEvent) object.LongpollBotResponse {
	return object.LongpollBotResponse{
		Ts:      strconv.FormatInt(time.Now().Unix(), 10),
		Updates: events,
		Failed:  0,
	}
}

func (s *subscription) returnEvents() (r object.LongpollBotResponse) {
	s.lock.Lock()
	defer s.lock.Unlock()

	r = pollResponse(s.events)
	s.events = nil
	return
}

func (s *subscription) Poll(ctxt context.Context) object.LongpollBotResponse {
	for {
		select {
		case _, ok := <-s.notify:
			if !ok {
				return pollResponse(nil)
			}

			return s.returnEvents()
		case <-ctxt.Done():
			return pollResponse(nil)
		}
	}
}

func (s *subscription) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	close(s.notify)
	return nil
}

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
